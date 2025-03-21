package main

import (
	"net/http"

	"github.com/dstrait38/receipt-processor/core"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var receipts = make(map[string]core.Receipt)

const notFound = "No receipt found for that ID."
const badRequest = "The receipt is invalid."

func main() {
	router := gin.Default()
	router.POST("/receipts/process", handleNewReceipt)
	router.GET("/receipts/:id/points", handleGetPoints)

	router.Run("localhost:8080")
}

func handleGetPoints(c *gin.Context) {
	id := c.Param("id")

	if receipt, exists := receipts[id]; exists {
		points := core.CalculatePoints(receipt)
		c.IndentedJSON(http.StatusOK, core.GetPointsResponse{Points: points})
	} else {
		c.IndentedJSON(http.StatusNotFound, core.ErrorResponse{Description: notFound})
	}

}

func handleNewReceipt(c *gin.Context) {
	var newReceipt core.Receipt

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatePrice", core.ValidatePrice)
	}

	if err := c.BindJSON(&newReceipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, core.ErrorResponse{Description: badRequest})
		return
	}

	newId := uuid.New().String()
	receipts[newId] = newReceipt
	response := core.NewReceiptResponse{Id: newId}

	c.IndentedJSON(http.StatusOK, response)
}
