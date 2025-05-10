package helper

const (
	algorithmLuhnCondition      = 9
	algorithmLuhnDecimalDivisor = 10
)

func ValidateAlgorithmLuhn(number uint64) bool {
	isEvenPosition := false
	var sum uint64 = 0

	for number > 0 {
		if isEvenPosition {
			value := number % algorithmLuhnDecimalDivisor
			value *= 2
			if value > algorithmLuhnCondition {
				value -= algorithmLuhnCondition
			}
			sum += value
			number /= algorithmLuhnDecimalDivisor
		} else {
			sum += number % algorithmLuhnDecimalDivisor
			number /= algorithmLuhnDecimalDivisor
		}
		isEvenPosition = !isEvenPosition
	}

	return sum%algorithmLuhnDecimalDivisor == 0
}
