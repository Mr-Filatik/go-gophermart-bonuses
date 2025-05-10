package helper

const (
	algorithmLuhnCondition      = 9
	algorithmLuhnDecimalDivisor = 10
)

func ValidateAlgorithmLuhn(number uint64) bool {
	isEvenPosition := false
	var sum uint64 = 0

	for number > 0 {
		value := number % algorithmLuhnDecimalDivisor
		if isEvenPosition {
			value *= 2
			if value > algorithmLuhnCondition {
				value -= algorithmLuhnCondition
			}
		}
		sum += value
		number /= algorithmLuhnDecimalDivisor
		isEvenPosition = !isEvenPosition
	}

	return sum%algorithmLuhnDecimalDivisor == 0
}
