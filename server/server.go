package server

import (
	"fmt"
	"net/http"
	"pdf_parser/parser"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllTransactions(c *gin.Context) {

	transactions, err := parser.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(transactions)
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func GetTransactionsByDateApi(c *gin.Context) {
	startDateParam := c.Query("startDate")
	if startDateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date parameter is required"})
		return
	}
	endDateParam := c.Query("endDate")
	if endDateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date parameter is required"})
		return
	}
	// fmt.Println(startDateParam, endDateParam)
	startDate, err := time.Parse("02-01-2006", startDateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date parameter is invalid"})
		return
	}

	endDate, err := time.Parse("02-01-2006", endDateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date parameter is invalid"})
		return
	}
	transactions, err := parser.GetTransactionsByDate(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(transactions)
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func GetBalanceByDateApi(c *gin.Context) {
	requiredDateParam := c.Query("date")
	if requiredDateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}
	requiredDate, err := time.Parse("02-01-2006", requiredDateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is invalid"})
		return
	}

	balance, err := parser.GetBalanceByDate(requiredDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(balance)
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/transactions", GetAllTransactions)
	router.GET("/transactions_by_date", GetTransactionsByDateApi)
	router.GET("/balance_by_date", GetBalanceByDateApi)
}
