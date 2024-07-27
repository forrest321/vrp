package calc

import (
	"fmt"
	d "github.com/forrest321/vrp/debug"
	t "github.com/forrest321/vrp/types"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Result struct {
	Solution []string
	Drivers  []Driver
}

type Driver struct {
	Loads      []Load
	TotalDist  float64
	CurrentPos t.Point
}

func (d *Driver) AcceptLoad(pickup t.Point, l Load) {
	d.TotalDist += Distance(pickup, l.Pickup) + Distance(l.Pickup, l.Dropoff)
	d.Loads = append(d.Loads, l)
	d.CurrentPos = l.Dropoff
}

func (d *Driver) FinalDistance() float64 {
	if len(d.Loads) == 0 {
		return 0
	}
	lastLoad := d.Loads[len(d.Loads)-1]
	return d.TotalDist + lastLoad.Dropoff.DistToDepot
}

type Load struct {
	Num     int
	Pickup  t.Point
	Dropoff t.Point
	Length  float64
}

var currentPosition = t.Depot

func Solve(remainingLoads []Load) Result {
	var drivers []Driver
	var acceptedLoads []Load
	for len(remainingLoads) > 0 {
		acceptedLoads = []Load{}
		currentDriver := Driver{CurrentPos: t.Depot}
		currentPosition = t.Depot
		loads := LoadsByPickup(remainingLoads)
		sort.Sort(loads)

		currentDriver.AcceptLoad(currentPosition, loads[0])
		currentPosition = loads[0].Dropoff
		remainingLoads = slices.Delete(loads, 0, 1)

		if len(remainingLoads) == 0 {
			drivers = append(drivers, currentDriver)
			break
		}
		lba := LoadsByCurrentPosition(remainingLoads)
		sort.Sort(lba)
		for _, l := range lba {
			pickupDist := Distance(currentPosition, l.Pickup)
			loadDist := Distance(l.Dropoff, l.Pickup)
			finalDist := Distance(t.Depot, l.Dropoff)
			totalDist := pickupDist + loadDist + finalDist + currentDriver.TotalDist

			d.Out(fmt.Sprintf("Distance if accepting next load: %v", totalDist))
			if totalDist > t.DriverMax {
				d.Out("rejecting load")
				continue
			}
			currentDriver.AcceptLoad(currentDriver.CurrentPos, l)
			currentPosition = l.Dropoff
			acceptedLoads = append(acceptedLoads, l)
		}
		drivers = append(drivers, currentDriver)
		for _, l := range acceptedLoads {
			i := slices.Index(remainingLoads, l)
			remainingLoads = slices.Delete(remainingLoads, i, i+1)
		}
	}

	var solution []string
	for _, dr := range drivers {
		innerArr := make([]string, len(dr.Loads))
		for i, l := range dr.Loads {
			innerArr[i] = strconv.Itoa(l.Num)
		}
		solLine := fmt.Sprintf("[%s]", strings.Join(innerArr, ", "))
		solution = append(solution, solLine)
	}

	return Result{Solution: solution, Drivers: drivers}
}
