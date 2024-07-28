package routing

import (
	"cmp"
	"fmt"
	c "github.com/forrest321/vrp/calc"
	t "github.com/forrest321/vrp/types"
	"slices"
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

	for len(remainingLoads) > 0 {
		acceptedLoads = []t.Load{}
		currentDriver = Driver{CurrentPos: t.Depot}
		//loads = t.LoadsByCurrentPosition(remainingLoads).SetCurrentPosition(currentDriver.CurrentPos)
		//sort.Sort(loads)
		sortFunc(remainingLoads, currentDriver.CurrentPos)

		currentDriver.AcceptLoad(remainingLoads[0])
		remainingLoads = slices.Delete(remainingLoads, 0, 1)

		sortFunc(remainingLoads, currentDriver.CurrentPos)

		for _, l := range remainingLoads {
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

func sortFunc(loads []t.Load, pos t.Point) func() {
	return func() {
		slices.SortFunc(loads, func(a, b t.Load) int {
			aPickupDist := c.Distance(a.Pickup.X, a.Pickup.Y, pos.X, pos.Y)
			bPickupDist := c.Distance(b.Pickup.X, b.Pickup.Y, pos.X, pos.Y)
			return cmp.Compare(aPickupDist, bPickupDist)
			//if n := cmp.Compare(aPickupDist, bPickupDist); n != 0 {
			//	return n
			//}
			//// If pickup distances are equal, order by dropoff distances
			//aDropoffDist := c.Distance(a.Dropoff.X, a.Dropoff.Y, pos.X, pos.Y)
			//bDropoffDist := c.Distance(b.Dropoff.X, b.Dropoff.Y, pos.X, pos.Y)
			//return cmp.Compare(aDropoffDist, bDropoffDist)
		})
	}
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
