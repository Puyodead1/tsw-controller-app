package math_utils

import "math"

func RoundToMarginOfError(value float64) float64 {
	return math.Round(value*10000.0) / 10000.0
}
