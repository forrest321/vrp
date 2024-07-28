package types

import "github.com/forrest321/vrp/calc"

var currentPosition = Depot

type LoadsByCurrentPosition []Load

func (ls LoadsByCurrentPosition) Len() int {
	return len(ls)
}

func (ls LoadsByCurrentPosition) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls LoadsByCurrentPosition) Less(i, j int) bool {
	iPickupDist := calc.Distance(ls[i].Pickup.X, ls[i].Pickup.Y, currentPosition.X, currentPosition.Y)
	jPickupDist := calc.Distance(ls[j].Pickup.X, ls[j].Pickup.Y, currentPosition.X, currentPosition.Y)
	return iPickupDist < jPickupDist
}

func (ls LoadsByCurrentPosition) SetCurrentPosition(p Point) LoadsByCurrentPosition {
	currentPosition = p
	return ls
}
