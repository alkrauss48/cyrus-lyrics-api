package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default()
	router.GET("/oauth/google", googleLogin)
	router.GET("/oauth/google/callback", googleLoginCallback)
	router.GET("/oauth/google/processed", googleLoginProcessed)
	router.GET("/sheets/new", newSheet)
	router.GET("/sheets/default", getDefaultSheetIds)
	router.GET("/sheets/:id", getSheetByID)

	router.Run("localhost:8000")
}
