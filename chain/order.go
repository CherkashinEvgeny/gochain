package chain

import "math"

type Order int

const (
	Impl      Order = math.MaxInt
	VeryLate  Order = 2
	Late      Order = 1
	Normal    Order = 0
	Early     Order = -1
	VeryEarly Order = -2
	Aspect    Order = math.MinInt
)
