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
parallel processing in this case would not work. The limiting factor, in this case, is that each driver must complete all 
loads before the end of their shift.

This solution iteratively builds a list of drivers and their routes by starting from the driver's origin, the Depot, finding 
the next nearest load, then finding the next nearest load to the previous load's dropoff location, and continuing in this fashion 
until there are no more loads that will fit into that driver's shift.

There are different data structures that would facilitate iteratively finding the next nearest load. In the interest of 
readability and simplicity, a simple slice is enough. Using Go's standard sort package provides a fast, highly optimized 
ability to sort a slice to find the next nearest load.
