package util

import (
	"crypto/rand"
	"mime/multipart"
	"murakali/config"

	"unicode"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
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

func UploadImageToCloudinary(c *gin.Context, cfg *config.Config, file multipart.File) string {
	cldService, _ := cloudinary.NewFromURL(cfg.External.CloudinaryURL)
	response, _ := cldService.Upload.Upload(c, file, uploader.UploadParams{})
	return response.SecureURL
}

func SKUGenerator(productName string) string {
	const otpChars = "1234567890"
	buffer := make([]byte, 8)
	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < 8; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return productName + "-" + string(buffer)
}
