package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	c "github.com/forrest321/vrp/calc"
	r "github.com/forrest321/vrp/routing"
	t "github.com/forrest321/vrp/types"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		exitWithError(t.UsageMessage, nil)
	}

	problemPath := flag.Args()[0]
	if ok, err := checkFile(problemPath); !ok || err != nil {
		exitWithError(t.NotFoundMessage, err)
	}

	loads, err := extractLoads(problemPath)
	if err != nil {
		exitWithError(t.DataErrorMessage, err)
	}

	sol := r.Solve(loads)

	for _, s := range sol {
		fmt.Println(s)
	}
}

func extractLoads(filePath string) ([]t.Load, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ' '
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if records[0][0] == t.FirstHeaderField {
		records = records[1:]
	}
	loads := make([]t.Load, len(records))
	for i, rr := range records {
		num, err := strconv.Atoi(rr[0])
		if err != nil {
			//no load num. data issue
			return nil, fmt.Errorf("%s %v, %w", t.ExtractLoadsErrorMessage, rr, err)
		}
		pickup, err := extractPoint(rr[1])
		if err != nil {
			return nil, err
		}
		dropoff, err := extractPoint(rr[2])
		if err != nil {
			return nil, err
		}
		length := c.Distance(dropoff.X, dropoff.Y, pickup.X, pickup.Y)
		loads[i] = t.Load{Num: num, Pickup: pickup, Dropoff: dropoff, Length: length}
	}
	return loads, nil
}

func extractPoint(input string) (t.Point, error) {
	var p t.Point
	input = strings.ReplaceAll(input, "(", "")
	input = strings.ReplaceAll(input, ")", "")
	xy := strings.Split(input, ",")
	x, err := strconv.ParseFloat(xy[0], 64)
	if err != nil {
		return p, fmt.Errorf("%s %v, %w", t.ExtractPointErrorMessage, xy[0], err)
	}
	p.X = x
	y, err := strconv.ParseFloat(xy[1], 64)
	if err != nil {
		return p, fmt.Errorf("%s %v, %w", t.ExtractPointErrorMessage, xy[1], err)
	}
	p.Y = y
	return p, nil
}

func checkFile(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func exitWithError(msg string, err error) {
	if len(msg) > 0 {
		fmt.Println(msg)
	}
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}
