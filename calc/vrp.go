package calc

import (
	"fmt"
	d "github.com/forrest321/vrp/debug"
	t "github.com/forrest321/vrp/types"
	"slices"
	"sort"
)

type driver struct {
	loads      []Load
	totalDist  float64
	currentPos t.Point
}

func (d *driver) AcceptLoad(l Load) {
	d.totalDist += l.Length + Distance(d.currentPos, l.Pickup)
	d.loads = append(d.loads, l)
	d.currentPos = l.Dropoff
}

func (d *driver) FinalDistance() float64 {
	if len(d.loads) == 0 {
		return 0
	}
	lastLoad := d.loads[len(d.loads)-1]
	return d.totalDist + lastLoad.Dropoff.DistToDepot
}

type Load struct {
	Num     int
	Pickup  t.Point
	Dropoff t.Point
	Length  float64
}

var currentPosition = t.Depot

func Solve(remainingLoads []Load) [][]int {
	var drivers []driver
	d.Out("Starting solve")
	d.Out("Number of loads: %d", len(remainingLoads))

	for len(remainingLoads) > 0 {
		d.Out("Number of remaining loads: %d", len(remainingLoads))
		currentDriver := driver{currentPos: t.Depot}
		d.Out("Number of Drivers: %d\n", len(drivers))
		loads := LoadsByPickup(remainingLoads)
		sort.Sort(loads)
		currentDriver.AcceptLoad(loads[0])
		currentPosition = loads[0].Dropoff
		remainingLoads = slices.Delete(loads, 0, 1)

		if len(remainingLoads) == 0 {
			break
		}

		for currentDriver.FinalDistance() < t.DriverMax {
			d.Out("Current Driver Distance: %.2f\n", currentDriver.FinalDistance())
			//lbpos := LoadsByCurrentPosition(remainingLoads)
			//sort.Sort(lbpos)
			//pickupDist := Distance(currentPosition, lbpos[0].Pickup)
			//d.Out(fmt.Sprintf("Distance if accepting next load: %v", currentDriver.totalDist+pickupDist+lbpos[0].Length+lbpos[0].Dropoff.DistToDepot))
			//if currentDriver.totalDist+pickupDist+lbpos[0].Length+lbpos[0].Dropoff.DistToDepot > t.DriverMax {
			//	d.Out("rejecting load")
			//	break
			//}
			//currentDriver.AcceptLoad(lbpos[0])
			//currentPosition = lbpos[0].Dropoff
			//remainingLoads = slices.Delete(lbpos, 0, 1)
			lba := LoadsByCurrentPosition(remainingLoads)
			sort.Sort(lba)
			pickupDist := Distance(currentPosition, lba[0].Pickup)
			d.Out(fmt.Sprintf("Distance if accepting next load: %v", currentDriver.totalDist+pickupDist+lba[0].Length+lba[0].Dropoff.DistToDepot))
			if currentDriver.totalDist+pickupDist+lba[0].Length+lba[0].Dropoff.DistToDepot > t.DriverMax {
				d.Out("rejecting load")
				break
			}
			currentDriver.AcceptLoad(lba[0])
			currentPosition = lba[0].Dropoff
			remainingLoads = slices.Delete(lba, 0, 1)
		}
		drivers = append(drivers, currentDriver)
	}

	var solution [][]int
	for j, dr := range drivers {
		d.Out("Driver #%v\n%+v\n", j, dr)
		innerArr := make([]int, len(dr.loads))
		for i, l := range dr.loads {
			d.Out(fmt.Sprintf("Load #%v\n%+v", i, l))
			innerArr[i] = l.Num
		}
		solution = append(solution, innerArr)
	}

	return solution
}
