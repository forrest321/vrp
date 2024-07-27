package calc

type LoadsByDepotPickup []Load

func (ls LoadsByDepotPickup) Len() int {
	return len(ls)
}

func (ls LoadsByDepotPickup) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls LoadsByDepotPickup) Less(i, j int) bool {
	return ls[i].Pickup.DistToDepot < ls[j].Pickup.DistToDepot
}

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

type LoadsByLength []Load

func (ls LoadsByLength) Len() int {
	return len(ls)
}

func (ls LoadsByLength) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls LoadsByLength) Less(i, j int) bool {
	return ls[i].Length < ls[j].Length
}

type LoadsByAll []Load

func (ls LoadsByAll) Len() int {
	return len(ls)
}

func (ls LoadsByAll) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls LoadsByAll) Less(i, j int) bool {
	var iTotal, jTotal float64
	iTotal = Distance(currentPosition, ls[i].Pickup) + ls[i].Length + ls[i].Dropoff.DistToDepot
	jTotal = Distance(currentPosition, ls[j].Pickup) + ls[j].Length + ls[j].Dropoff.DistToDepot
	return iTotal < jTotal
}

type LoadsByPickup []Load

func (ls LoadsByPickup) Len() int {
	return len(ls)
}

func (ls LoadsByPickup) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls LoadsByPickup) Less(i, j int) bool {
	var iTotal, jTotal float64
	iTotal = ls[i].Pickup.DistToDepot + ls[i].Length
	jTotal = ls[j].Pickup.DistToDepot + ls[j].Length
	return iTotal < jTotal
}
