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
	return d.TotalDist + Distance(d.Loads[len(d.Loads)-1].Dropoff, t.Depot)
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
	var currentDriver Driver
	var pickupDist, finalDist, totalDist float64
	var loads, lba LoadsByCurrentPosition
	for len(remainingLoads) > 0 {
		acceptedLoads = []Load{}
		currentDriver = Driver{CurrentPos: t.Depot}
		currentPosition = t.Depot
		loads = remainingLoads
		sort.Sort(loads)

		acceptLoad(&currentDriver, loads[0])
		remainingLoads = slices.Delete(loads, 0, 1)

		if len(remainingLoads) == 0 {
			//could continue through and exit normally, but why wait?
			drivers = append(drivers, currentDriver)
			break
		}
		lba = remainingLoads
		sort.Sort(lba)
		for _, l := range lba {
			pickupDist = Distance(currentPosition, l.Pickup)
			finalDist = Distance(t.Depot, l.Dropoff)
			totalDist = pickupDist + l.Length + finalDist + currentDriver.TotalDist

			d.Out(fmt.Sprintf("Distance if accepting next load: %v", totalDist))
			if totalDist > t.DriverMax {
				d.Out("rejecting load")
				continue
			}
			acceptLoad(&currentDriver, l)
			acceptedLoads = append(acceptedLoads, l)
		}
		drivers = append(drivers, currentDriver)
		for _, l := range acceptedLoads {
			//remove loads accepted in range above
			i := slices.Index(remainingLoads, l)
			if i == -1 {
				//shouldn't happen. if it does, prob panic
				continue
			}
			remainingLoads = slices.Delete(remainingLoads, i, i+1)
		}
	}

	return Result{Solution: formatSolution(drivers), Drivers: drivers}
}

func acceptLoad(dr *Driver, l Load) {
	dr.AcceptLoad(currentPosition, l)
	currentPosition = l.Dropoff
}

func formatSolution(drivers []Driver) []string {
	var solution []string
	for _, dr := range drivers {
		innerArr := make([]string, len(dr.Loads))
		for i, l := range dr.Loads {
			innerArr[i] = strconv.Itoa(l.Num)
		}
		solLine := fmt.Sprintf("[%s]", strings.Join(innerArr, ", "))
		solution = append(solution, solLine)
	}
	return solution
}
