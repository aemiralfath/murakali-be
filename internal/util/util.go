package util

import (
	"crypto/rand"
	"unicode"
)

func GenerateOTP(length int) (string, error) {
	const otpChars = "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func VerifyPassword(s string) bool {
	if len(s) < 8 || len(s) > 40 {
		return false
	}

	number, upper, special := false, false, false

	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case isSpecialCharacter(c):
			special = true
		case number && upper && special:
			break
		}
	}

	return number && upper && special
}

func isSpecialCharacter(char rune) bool {
	return !unicode.IsLetter(char) && !unicode.IsNumber(char) && !unicode.IsSpace(char)
}
