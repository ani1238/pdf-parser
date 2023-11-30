package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pdf_parser/emailfetcher"
	"pdf_parser/parser"
	"pdf_parser/server"

	"github.com/gin-gonic/gin"
)

func main() {
	userEmail := "anisumi1238@gmail.com"
	subject := "Bank Statement Attached Valyx"
	err := emailfetcher.FetchPdfsFromEmailForSubject(userEmail, subject)
	if err != nil {
		panic(err)
	}

	folderPath := "./outpdfs"

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
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
		panic(err)
	}

	router := gin.Default()

	server.SetupRoutes(router)

	fmt.Println("Server is running on http://localhost:8080")
	router.Run(":8080")

	return
}
