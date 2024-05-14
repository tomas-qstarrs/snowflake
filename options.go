package snowflake

const (
	WaitMethodSleep = iota
	WaitMethodSpin
)

type Options struct {
	pattern     *Pattern
	safe        bool
	waitMethod  int
	etcd        *Etcd
	node        int64
	networkNode bool
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

var defaultOptions = Options{
	pattern:     DefaultPattern,
	safe:        true,
	waitMethod:  WaitMethodSleep,
	etcd:        nil,
	node:        0,
	networkNode: false,
}

type Option func(*Options)

func WithPattern(p *Pattern) Option {
	return func(o *Options) {
		o.pattern = p
	}
}

func WithSafe(safe bool) Option {
	return func(o *Options) {
		o.safe = safe
	}
}

func WithWaitMethod(method int) Option {
	return func(o *Options) {
		o.waitMethod = method
	}
}

func WithEtcdNode(etcd *Etcd) Option {
	return func(o *Options) {
		o.etcd = etcd
	}
}

func WithNetworkNode() Option {
	return func(o *Options) {
		o.networkNode = true
	}
}

func WithNode(node int64) Option {
	return func(o *Options) {
		o.node = node
	}
}
