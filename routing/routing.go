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

func (d *Driver) CanAcceptLoad(l t.Load) bool {
	pickupDist := distance(d.CurrentPos, l.Pickup)
	dropoffDist := distance(l.Dropoff, t.Depot)
	if pickupDist+l.Length+dropoffDist+d.TotalDist <= t.DriverMax {
		return true
	}
	return false
}

func Solve(remainingLoads []t.Load) []string {
	var drivers []*Driver
	var currentDriver *Driver
	var driverIsDone = false
	var lbp, acceptedLoads t.LoadsByCurrentPosition

	lbp = remainingLoads

	for len(lbp) > 0 {
		currentDriver = &Driver{CurrentPos: t.Depot}
		drivers = append(drivers, currentDriver)
		driverIsDone = false

		for !driverIsDone {
			lbp = lbp.SetCurrentPosition(currentDriver.CurrentPos)
			sort.Sort(lbp)
			acceptedLoads = []t.Load{}

			for i, l := range lbp {
				if currentDriver.CanAcceptLoad(l) {
					currentDriver.AcceptLoad(l)
					acceptedLoads = append(acceptedLoads, l)
					break
				}
				if i == len(lbp)-1 {
					driverIsDone = true
				}
			}
			lbp = removeLoads(lbp, acceptedLoads)
			if len(lbp) == 0 {
				driverIsDone = true
			}
		}
	}

	return formatSolution(drivers)
}

func removeLoads(loads, toRemove []t.Load) []t.Load {
	var remainingLoads []t.Load
	for _, load := range loads {
		if !slices.Contains(toRemove, load) {
			remainingLoads = append(remainingLoads, load)
		}
	}
	return remainingLoads
}

func formatSolution(drivers []*Driver) []string {
	var solution, innerArr []string
	var solLine string
	for _, dr := range drivers {
		innerArr = make([]string, len(dr.Loads))
		for i, l := range dr.Loads {
			innerArr[i] = strconv.Itoa(l.Num)
		}
		solLine = fmt.Sprintf("[%s]", strings.Join(innerArr, ", "))
		solution = append(solution, solLine)
	}
	return solution
}

func distance(p1, p2 t.Point) float64 {
	return c.Distance(p1.X, p1.Y, p2.X, p2.Y)
}
