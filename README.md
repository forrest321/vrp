# VRP
Vehicle Routing Problem - a simple solution with 0 external dependencies

## Problem Description
> For this challenge, you will submit a program that solves a version of the Vehicle Routing Problem (VRP).
> 
> The VRP specifies a set of loads to be completed efficiently by an unbounded number of drivers.
> 
> Each load has a pickup location and a dropoff location, each specified by a Cartesian point. A driver completes a load by
> driving to the pickup location, picking up the load, driving to the dropoff, and dropping off the load. The time required 
> to drive from one point to another, in minutes, is the Euclidean distance between them.
> 
> Each driver starts and ends his shift at a depot located at (0,0). A driver may complete multiple loads on his shift, but 
> may not exceed 12 hours of total drive time. That is, the total Euclidean distance of completing all his loads, including 
> the return to (0,0), must be less than 12*60.
> 
> A VRP solution contains a list of drivers, each of which has an ordered list of loads to be completed. All loads must be 
> assigned to a driver.

## Approach
[The simplest approach to solve VRP is to use the Nearest Neighbor Method](https://arxiv.org/pdf/2303.04147).  This solution 
relies on finding the next nearest route for a driver, then the next, and so on. To further optimize this, given a fixed number 
of drivers, the routes can be found in parallel. Since the problem description specifies the number of drivers is unbounded,
parallel processing in this case would not work. The limiting factor for this problem is that each driver must complete all 
loads before the end of their shift.

This solution iteratively builds a list of drivers and their routes by starting from the driver's origin, the Depot, finding 
the nearest load, then finding the next nearest load to the previous load's dropoff location, and continuing in this fashion 
until there are no more loads that will fit into that driver's shift. At this point, if there are any remaining loads, the 
process repeats with a new driver.

There are different data structures that would facilitate iteratively finding the next nearest load. In the interest of 
readability and simplicity, a simple slice is enough. Using Go's standard sort package provides a fast, highly optimized 
ability to sort a slice to find the next nearest load.

## sort.Sort
By implementing [sort.Interface](https://pkg.go.dev/sort#pkg-overview), use of the sort.Sort package is trivial. The internal sort algorithm is optimized and outperforms 
other popular sorting algo's, in some cases by order(s) of magnitude. During testing of this solution, sorting by many 
different parameters was tested, such as: start location, load distance, distance to pickup + load distance, and round trip 
distance. Sorting by the next closest load consistently provided the best mean cost and run time.

## sort.Sort vs slices.SortFunc
According to [Go documentation](https://pkg.go.dev/sort#Sort), "in many situations, the newer slices.SortFunc 
function is more ergonomic and runs faster [than sort.Sort]."  This did happen to be the case, however the resulting increase 
of mean cost by $14,799.10 did not justify the decrease in run time by 0.69ms.

|  Sort Method | Mean Cost | Mean Run Time | 
|--------------| --- | --- |
|sort.Func | 70757.57101574977 | 1.2691617012023926ms |
|sort.Sort | 55958.475762225324 | 1.9564032554626465ms |

# Usage
## Download the source:
To download a zip [click here](https://github.com/forrest321/vrp/archive/refs/heads/main.zip), or you can clone the repo:
```shell
git clone https://github.com/forrest321/vrp.git
```
## Building
After cloning the repo, cd into the folder and run `go build`:
```shell
cd vrp
go build
```
An executable named `vrp` or `vrp.exe` will be created. It can be copied anywhere and called by `./vrp` in *nix based systems 
or `vrp` on Windows. On *nix systems, copying to a location within the `$PATH` allows it to be called without the leading 
`./`, so just `vrp`.
## Installation
If you just want an executable, go to [the releases page](https://github.com/forrest321/vrp/releases) and download the appropriate 
file for your system.  

If you have Go installed, you can install vrp directly via:
```shell
go install github.com/forrest321/vrp
```
and run it with:
```shell
vrp /path/to/problemFile.txt
```

A list of sample problem files in the expected format is available in the `/problems` folder in this repo. If you have cloned the repo, 
the entire process will look like this:
```shell
git clone https://github.com/forrest321/vrp.git
cd vrp
go build
./vrp problems/problem1.txt
```
If everything worked, you should see the outputted routes:
```shell
[1, 4, 10]
[8, 2, 5]
[7]
[6]
[3]
[9]
```
