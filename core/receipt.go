package core

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Receipt struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required" time_format:"2006-01-02"`
	PurchaseTime string `json:"purchaseTime" binding:"required" time_format:"15:04"`
	Items        []Item `json:"items" binding:"required,min=1,dive"`
	Total        string `json:"total" binding:"required,validatePrice"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required,validatePrice"`
}

type NewReceiptResponse struct {
	Id string `json:"id"`
}

type GetPointsResponse struct {
	Points int `json:"points"`
}

type ErrorResponse struct {
	Description string `json:"description"`
}

// Custom validator for fields that require two decimal places
func ValidatePrice(fl validator.FieldLevel) bool {
	price := fl.Field().String()
	re := regexp.MustCompile(`^\d+\.\d{2}$`)
	return re.MatchString(price)
}

func CalculatePoints(receipt Receipt) int {
	var points = 0
	retailerSplit := strings.SplitSeq(receipt.Retailer, " ")
	for word := range retailerSplit {
		points += totalAlphaNum(word)
	}

	total, _ := strconv.ParseFloat(receipt.Total, 64)

	if total == math.Floor(total) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		if len(strings.Trim(item.ShortDescription, " "))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	pdLen := len(receipt.PurchaseDate)
	day, _ := strconv.ParseInt(receipt.PurchaseDate[pdLen-1:], 10, 64)

	if day%2 == 1 {
		points += 6
	}

	time, _ := strconv.ParseInt(receipt.PurchaseTime[:2], 10, 64)

	if time-14 >= 0 && time-14 <= 2 {
		points += 10
	}

	return points
}

func totalAlphaNum(str string) int {
	var cnt = 0
	for _, char := range str {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			cnt++
		}
	}
	return cnt
}
