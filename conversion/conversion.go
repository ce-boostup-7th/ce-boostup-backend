package conversion

import (
	"strconv"
)

// StringToInt convert type string to int
func StringToInt(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// StringToFloat convert type string to float64
func StringToFloat(str string) (float64) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return value
}
