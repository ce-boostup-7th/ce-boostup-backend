package conversion

import (
	"log"
	"strconv"
)

// StringToInt convert type string to int
func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

// StringToFloat convert type string to float64
func StringToFloat(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
		return 0.0
	}
	return value
}
