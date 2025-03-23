package core

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Receipt struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required,validateDate"`
	PurchaseTime string `json:"purchaseTime" binding:"required,validateTime"`
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

// Custom validator for PurcahseDate field
func ValidateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	validString := re.MatchString(date)
	_, er := time.Parse("2006-01-02", date)
	return validString && er == nil
}

// Custom validator for PurchaseTime field
func ValidateTime(fl validator.FieldLevel) bool {
	timeField := fl.Field().String()
	re := regexp.MustCompile(`^\d{2}:\d{2}$`)
	validString := re.MatchString(timeField)
	_, er := time.Parse("15:04", timeField)
	return validString && er == nil
}

// Function to calcualte points for receipt per challenge requirements
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
