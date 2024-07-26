package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x, y float64
}
type load struct {
	num     int
	pickup  point
	dropoff point
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: vrp [filename.txt]")
		os.Exit(1)
	}
	probPath := flag.Args()[0]
	if ok, err := checkFile(probPath); !ok || err != nil {
		fmt.Println("File not found")
		os.Exit(1)
	}
	fmt.Println("Loading problem file: ", probPath)
	loads, err := extractLoads(probPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Loads: \n%+v\n", loads)
}

func extractLoads(filePath string) ([]load, error) {
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
	loads := make([]load, len(records))
	for i, r := range records {
		if i == 0 {
			continue
		}
		num, err := strconv.Atoi(r[0])
		if err != nil {
			//how to handle?
			continue
		}
		pickup := extractPoint(r[1])
		dropoff := extractPoint(r[2])
		loads[i] = load{num: num, pickup: pickup, dropoff: dropoff}
	}

	return loads, nil
}

func extractPoint(input string) point {
	var p point
	input = strings.ReplaceAll(input, "(", "")
	input = strings.ReplaceAll(input, ")", "")
	xy := strings.Split(input, ",")
	x, err := strconv.ParseFloat(xy[0], 64)
	if err != nil {
		//how to handle?
		return p
	}
	p.x = x
	y, err := strconv.ParseFloat(xy[1], 64)
	if err != nil {
		//how to handle?
		return p
	}
	p.y = y
	return p
}

func checkFile(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
