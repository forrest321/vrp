package types

type Point struct {
	X, Y float64
}

type Result struct {
	Solution []string
	Drivers  []Driver
}

type Driver struct {
	Loads      []Load
	TotalDist  float64
	CurrentPos Point
}

type Load struct {
	Num     int
	Pickup  Point
	Dropoff Point
	Length  float64
}
