package calc

import (
	"math"
)

func Distance(x1, y1, x2, y2 float64) float64 {
	//(x1, y1) to (x2, y2) takes sqrt((x2-x1)^2 + (y2-y1)^2)
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}
