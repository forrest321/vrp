package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	c "github.com/forrest321/vrp/calc"
	d "github.com/forrest321/vrp/debug"
	t "github.com/forrest321/vrp/types"
	"os"
	"strconv"
	"strings"
)

var verbose bool = false

func main() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
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
	d.Out(fmt.Sprintf("Loading problem file: %s", probPath))

	loads, err := extractLoads(probPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	d.Out(fmt.Sprintf("Loads: \n%+v", loads))
	sol := c.Solve(loads)
	for _, s := range sol {
		fmt.Println(s)
	}
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

		d.Out(fmt.Sprintf("load[%d] = %+v", i, loads[i]))
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
	p.DistToDepot = c.Distance(t.Depot, p)
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
