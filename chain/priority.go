package chain

import "math"

type Priority int

const (
	Impl    Priority = math.MaxInt
	Highest Priority = 2
	Height  Priority = 1
	Medium  Priority = 0
	Low     Priority = -1
	Lowest  Priority = -2
	Aspect  Priority = math.MinInt
)
