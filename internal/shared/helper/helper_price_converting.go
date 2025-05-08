package helper

import "math"

const mantissaSize = 100

func ConvertPriceToFloat64(price uint64) float64 {
	return float64(price) / mantissaSize
}

func ConvertPriceToUint64(price float64) uint64 {
	return uint64(math.Round(price * mantissaSize))
}
