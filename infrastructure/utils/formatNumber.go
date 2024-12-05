package utils

import "strconv"

func FormatNumber(number int) string {
	formattedNumber := strconv.FormatInt(int64(number), 10)
	for i := len(formattedNumber) - 3; i > 0; i -= 3 {
		formattedNumber = formattedNumber[:i] + "," + formattedNumber[i:]
	}
	return formattedNumber
}
