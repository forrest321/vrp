package types

var Depot = Point{X: 0, Y: 0}

const DriverMax = float64(60 * 12)
const UsageMessage = "Usage: vrp /path/to/filename.txt"
const NotFoundMessage = "File not found"
const DataErrorMessage = "Data error"
const ExtractLoadsErrorMessage = "extractLoads error parsing record:"
const ExtractPointErrorMessage = "extractPoint error parsing point:"

const FirstHeaderField = "loadNumber"
