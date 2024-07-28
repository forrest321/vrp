package routing

import (
	"fmt"
	c "github.com/forrest321/vrp/calc"
	t "github.com/forrest321/vrp/types"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Driver t.Driver

func (d *Driver) AcceptLoad(l t.Load) {
	d.TotalDist += distance(d.CurrentPos, l.Pickup) + distance(l.Pickup, l.Dropoff)
	d.Loads = append(d.Loads, l)
	d.CurrentPos = l.Dropoff
}

func (d *Driver) FinalDistance() float64 {
	if len(d.Loads) == 0 {
		return 0
	}
	return d.TotalDist + distance(d.Loads[len(d.Loads)-1].Dropoff, t.Depot)
}

func Solve(remainingLoads []t.Load) []string {
	var drivers []Driver
	var acceptedLoads []t.Load
	var currentDriver Driver
	var pickupDist, finalDist, totalDist float64
	var loads, lba t.LoadsByCurrentPosition

	for len(remainingLoads) > 0 {
		acceptedLoads = []t.Load{}
		currentDriver = Driver{CurrentPos: t.Depot}
		loads = t.LoadsByCurrentPosition(remainingLoads).SetCurrentPosition(currentDriver.CurrentPos)
		sort.Sort(loads)

		currentDriver.AcceptLoad(loads[0])
		remainingLoads = slices.Delete(loads, 0, 1)

		lba = t.LoadsByCurrentPosition(remainingLoads).SetCurrentPosition(currentDriver.CurrentPos)
		sort.Sort(lba)
		for _, l := range lba {
			pickupDist = distance(currentDriver.CurrentPos, l.Pickup)
			finalDist = distance(t.Depot, l.Dropoff)
			totalDist = pickupDist + l.Length + finalDist + currentDriver.TotalDist

			if totalDist > t.DriverMax {
				continue
			}
			currentDriver.AcceptLoad(l)
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

	return formatSolution(drivers)
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

func distance(p1, p2 t.Point) float64 {
	return c.Distance(p1.X, p1.Y, p2.X, p2.Y)
}
