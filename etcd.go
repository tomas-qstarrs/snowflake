package snowflake

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	Prefix        string
	Addr          string
	Timeout       time.Duration
	LeaseTime     time.Duration
	LeaseInterval time.Duration
}

func (n *Node) generateNodeFromEtcd() {
	var (
		prefix        = n.etcd.Prefix
		addr          = n.etcd.Addr
		timeout       = n.etcd.Timeout
		leaseTime     = n.etcd.LeaseTime
		leaseInterval = n.etcd.LeaseInterval
	)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := doWithContext(ctx, func() error {
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   strings.Split(addr, ","),
				DialTimeout: timeout,
			})
			if err != nil {
				return err
			}
			defer client.Close()

			resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
			if err != nil {
				return err
			}

			existNodeMap := make(map[int]int)
			for _, ev := range resp.Kvs {
				num, _ := strconv.Atoi(string(ev.Value))
				existNodeMap[num] = num
			}
			for id := 0; id <= int(n.nodeMax); id++ {
				if _, ok := existNodeMap[id]; !ok {
					n.node = int64(id)
					break
				}
			}
			if n.node == n.nodeMax {
				return fmt.Errorf("no available node id")
			}

			lease := clientv3.NewLease(client)
			grant, err := lease.Grant(ctx, int64(leaseTime.Seconds()))
			if err != nil {
				return err
			}
			_, err = client.Put(ctx, prefix+strconv.FormatInt(n.node, 10), strconv.FormatInt(n.node, 10), clientv3.WithLease(grant.ID))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Println(err)
		} else {
			break
		}
	}

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if err := doWithContext(ctx, func() error {
				client, err := clientv3.New(clientv3.Config{
					Endpoints:   strings.Split(addr, ","),
					DialTimeout: timeout,
				})
				if err != nil {
					return err
				}
				defer client.Close()

				lease := clientv3.NewLease(client)
				grant, err := lease.Grant(ctx, int64(leaseTime.Seconds()))
				if err != nil {
					return err
				}
				_, err = client.Put(ctx, prefix+strconv.FormatInt(n.node, 10), strconv.FormatInt(n.node, 10), clientv3.WithLease(grant.ID))
				if err != nil {
					return err
				}
				time.Sleep(time.Second * time.Duration(leaseInterval))
				return nil
			}); err != nil {
				log.Println(err)
			}
		}
	}()
}

func do(a any) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panic: %v", v)
		}
	}()

	switch f := a.(type) {
	case func():
		f()
	case func() error:
		err = f()
	default:
		panic(fmt.Errorf("invalid function type: %T", f))
	}

	return
}

func doWithContext(ctx context.Context, a any) (err error) {
	errCh := make(chan error, 1)

	go func() {
		errCh <- do(a)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errCh:
		return err
	}
}

func doWithTimeout(d time.Duration, a any) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	return doWithContext(ctx, a)
}
