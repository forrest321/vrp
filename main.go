package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	c "github.com/forrest321/vrp/calc"
	d "github.com/forrest321/vrp/debug"
	t "github.com/forrest321/vrp/types"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var verbose bool = false
var loadIds []int

func main() {
	check := false
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&check, "c", false, "check solution")
	flag.Parse()
	d.SetDebug(verbose)
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: vrp [filename.txt]")
		os.Exit(1)
	}
	probPath := flag.Args()[0]
	if ok, err := checkFile(probPath); !ok || err != nil {
		fmt.Println("File not found")
		os.Exit(1)
	}

	loads, err := extractLoads(probPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sol := c.Solve(loads)
	if check {
		checkResult(sol)
		os.Exit(0)
	}

	for _, s := range sol.Solution {
		fmt.Println(s)
	}
}

func checkResult(r c.Result) {
	checkSolution(r.Solution)
	checkDrivers(r.Drivers)
}

func checkDrivers(d []c.Driver) {
	fmt.Printf("Number of drivers: %v\n", len(d))
	//find driver's driving time
	for j, dr := range d {
		fmt.Printf("Checking driver# %v\n", j)
		var routes []int
		dist := float64(0)
		for i, l := range dr.Loads {
			routes = append(routes, l.Num)
			if i == 0 {
				dist += c.Distance(t.Depot, l.Pickup)
			}
			dist += c.Distance(l.Dropoff, l.Pickup)
			if i == len(dr.Loads)-1 {
				dist += c.Distance(t.Depot, l.Dropoff)
			}
		}
		fmt.Printf("Routes: %v\n", routes)
		fmt.Printf("Reported driver distance: %v\n", dr.FinalDistance())
		fmt.Printf("Actual Driver Distance: %v is over limit? %v\n", dist, dist > t.DriverMax)
		fmt.Printf("Distances match: %v\n", dist == dr.FinalDistance())
	}
}

func checkSolution(s []string) {
	fmt.Println("checking solution")
	fmt.Printf("number of loads: %v\n", len(s))
	var solIds []int
	for _, line := range s {
		line = strings.ReplaceAll(line, "[", "")
		line = strings.ReplaceAll(line, "]", "")
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\n", "")
		ids := strings.Split(line, ",")
		for _, id := range ids {
			i, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
			}
			solIds = append(solIds, i)
		}
	}
	sort.Ints(solIds)
	sort.Ints(loadIds)
	missingFromSolution := missingNumbers(solIds)
	fmt.Printf("missing from solution ids: %v\n", missingFromSolution)
	missingFromProblem := missingNumbers(loadIds)
	fmt.Printf("missing from problem ids: %v\n", missingFromProblem)
	if slices.Equal(solIds, loadIds) {
		fmt.Println("Solutions are equal")
		return
	}
	var missingIds []int
	for _, id := range loadIds {
		if slices.Index(solIds, id) == -1 {
			missingIds = append(missingIds, id)
		}
	}
	if len(missingIds) > 0 {
		fmt.Printf("IDS are missing: %v\n", missingIds)
		os.Exit(1)
	}
}

func missingNumbers(x []int) []int {
	if len(x) == 0 {
		return nil
	}
	sort.Ints(x)
	var missing []int
	//assume range is x[0] through x[len(x)-1]
	for i := x[0]; i <= x[len(x)-1]; i++ {
		if _, ok := slices.BinarySearch(x, i); !ok {
			missing = append(missing, i)
		}
	}

	return missing
}

func extractLoads(filePath string) ([]c.Load, error) {
	file, err := os.Open(filePath)
	if err != nil {
		d.Out(fmt.Sprintf("Error opening file: %s", err))
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ' '
	records, err := csvReader.ReadAll()
	if err != nil {
		d.Out(fmt.Sprintf("Error reading csv: %s", err))
		return nil, err
	}
	if records[0][0] == "loadNumber" {
		records = records[1:]
	}
	loads := make([]c.Load, len(records))
	for i, r := range records {
		num, err := strconv.Atoi(r[0])
		if err != nil {
			//how to handle?
			continue
		}
		pickup := extractPoint(r[1])
		dropoff := extractPoint(r[2])
		length := c.Distance(dropoff, pickup)
		loads[i] = c.Load{Num: num, Pickup: pickup, Dropoff: dropoff, Length: length}
		loadIds = append(loadIds, num)
	}
	return loads, nil
}

func extractPoint(input string) t.Point {
	var p t.Point
	input = strings.ReplaceAll(input, "(", "")
	input = strings.ReplaceAll(input, ")", "")
	xy := strings.Split(input, ",")
	x, err := strconv.ParseFloat(xy[0], 64)
	if err != nil {
		//how to handle?
		if verbose {
			fmt.Println("Error converting input to float: ", err)
		}
		return p
	}
	p.X = x
	y, err := strconv.ParseFloat(xy[1], 64)
	if err != nil {
		//how to handle?
		if verbose {
			fmt.Println("Error converting input to float: ", err)
		}
		return p
	}
	p.Y = y
	return p
}

func checkFile(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if verbose {
			fmt.Println("File not found: ", filename)
		}
		return false, err
	}
	return true, nil
}
