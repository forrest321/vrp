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

	lbp := t.LoadsByCurrentPosition(remainingLoads).SetCurrentPosition(t.Depot)
	sort.Sort(lbp)
	for len(lbp) > 0 {
		acceptedLoads = []t.Load{}
		currentDriver = Driver{CurrentPos: t.Depot}

		lbp = lbp.SetCurrentPosition(currentDriver.CurrentPos)
		sort.Sort(lbp)

		currentDriver.AcceptLoad(lbp[0])
		lbp = slices.Delete(lbp, 0, 1)

		lbp = lbp.SetCurrentPosition(currentDriver.CurrentPos)
		sort.Sort(lbp)
		for _, l := range lbp {
			pickupDist = distance(currentDriver.CurrentPos, l.Pickup)
			finalDist = distance(l.Dropoff, t.Depot)
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
			i := slices.Index(lbp, l)
			if i == -1 {
				//shouldn't happen. if it does, prob panic
				continue
			}
			lbp = slices.Delete(lbp, i, i+1)
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
