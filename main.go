package main

import (
	"fmt"
	"pdf_parser/emailfetcher"
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

	router := gin.Default()

	server.SetupRoutes(router)

	fmt.Println("Server is running on http://localhost:8080")
	router.Run(":8080")

	return
}
