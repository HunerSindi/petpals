
package utils

import "strconv"

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0 // or handle error appropriately
	}
	return f
}
