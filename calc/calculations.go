package calc

import (
	t "github.com/forrest321/vrp/types"
	"math"
)

func Distance(p1, p2 t.Point) float64 {
	//(x1, y1) to (x2, y2) takes sqrt((x2-x1)^2 + (y2-y1)^2)
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}
