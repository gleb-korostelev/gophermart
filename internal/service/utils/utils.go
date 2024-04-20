package utils

func IsValidOrderID(id string) bool {
	var sum int
	nDigits := len(id)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit := int(id[i] - '0')
		if i%2 == parity {
			digit = digit * 2
		}
		if digit > 9 {
			digit -= 9
		}
		sum += digit
	}
	return sum%10 == 0
}
