package calc

import (
	t "github.com/forrest321/vrp/types"
	"math"
)

func Distance(p1, p2 t.Point) float64 {
	//(x1, y1) to (x2, y2) takes sqrt((x2-x1)^2 + (y2-y1)^2)
	return math.Sqrt(((p2.X - p1.X) * (p2.X - p1.X)) + ((p2.Y - p1.Y) * (p2.Y - p1.Y)))
}

//func TotalCost(n, miles float64) float64 {
//	//total_cost = 500*number_of_drivers + total_number_of_driven_minutes
//	return (500 * n) + miles
//}
