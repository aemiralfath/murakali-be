package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"mime/multipart"
	"murakali/config"
	"murakali/internal/model"
	"time"

	"unicode"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/sony/sonyflake"
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

func GenerateRandomAlpaNumeric(length int) (string, error) {
	const alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	charsLength := len(alphaNum)
	for i := 0; i < length; i++ {
		buffer[i] = alphaNum[int(buffer[i])%charsLength]
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
func GenerateInvoice() (string, error) {

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return "", err
	}
	invoice := fmt.Sprintf("INV/%s/MRK/%x\n", time.Now().Format("20060102"), id)

	return invoice, nil
}

func CalculateDiscount(price float64, disc *model.Discount) (float64, float64) {
	if disc.MaxDiscountPrice == nil {
		return 0, price
	}
	var maxDiscountPrice float64
	var minProductPrice float64
	var discountPercentage float64
	var discountFixPrice float64

	var resultDiscount float64

	maxDiscountPrice = *disc.MaxDiscountPrice
	if maxDiscountPrice == 0 {
		return 0, price
	}

	if disc.MinProductPrice != nil {
		minProductPrice = *disc.MinProductPrice
	}

	if disc.DiscountPercentage != nil {
		discountPercentage = *disc.DiscountPercentage
		if price >= minProductPrice && discountPercentage > 0 {
			resultDiscount = math.Min(maxDiscountPrice, price * (discountPercentage/100.00))
		}
	}

	if disc.DiscountFixPrice != nil {
		discountFixPrice = *disc.DiscountFixPrice
		if price >= minProductPrice && discountFixPrice > 0 {
			resultDiscount = math.Max(resultDiscount, discountFixPrice)
			resultDiscount = math.Min(resultDiscount, maxDiscountPrice)
		}
	}
	return resultDiscount, price - resultDiscount
}