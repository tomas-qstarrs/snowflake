package snowflake

import "time"

// Pattern is a snowflake pattern
type Pattern struct {
	Epoch    time.Time
	Tick     time.Duration
	StepBits [2]uint8
	NodeBits [2]uint8
	TimeBits [2]uint8
}

var (
	DefaultPattern = NewPattern(
		time.Unix(1288834974, 657000000),
		time.Millisecond,
		[2]uint8{0, 12},  // 12 bits
		[2]uint8{12, 10}, // 10 bits
		[2]uint8{22, 41}, // 41 bits
	)
)

// NewPattern returns a new snowflake pattern
func NewPattern(epoch time.Time, tick time.Duration, stepBits, nodeBits, timeBits [2]uint8) *Pattern {
	return &Pattern{
		Epoch:    epoch,
		Tick:     tick,
		StepBits: stepBits,
		TimeBits: timeBits,
		NodeBits: nodeBits,
	}
}

// NewNode returns a new snowflake node that can be used to generate snowflake
func (p *Pattern) NewNode(opts ...Option) (*Node, error) {
	return NewNode(append(opts, WithPattern(p))...)
}
