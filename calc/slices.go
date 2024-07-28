package calc

type LoadsByCurrentPosition []Load

func (ls LoadsByCurrentPosition) Len() int {
	return len(ls)
}

func (ls LoadsByCurrentPosition) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls LoadsByCurrentPosition) Less(i, j int) bool {
	iPickupDist := Distance(ls[i].Pickup, currentPosition)
	jPickupDist := Distance(ls[j].Pickup, currentPosition)
	return iPickupDist < jPickupDist
}
