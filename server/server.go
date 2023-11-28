package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"pdf_parser/parser"

	"github.com/gin-gonic/gin"
)

func GetAllTransactions(c *gin.Context) {
	folderPath := "./outpdfs"

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		err = parser.ParsePdf(path)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the folder:", err)
	}

	transactions, err := parser.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(transactions)
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/transactions", GetAllTransactions)
	// router.GET("/active_players", GetActivePlayers)
}
