package main

import (
	"log"

	"github.com/alkrauss48/cyrus-lyrics-api/oauth"
	"github.com/alkrauss48/cyrus-lyrics-api/public"
	"github.com/alkrauss48/cyrus-lyrics-api/sheets"
	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	// Public, unauthenticated routes
	router.GET("/", public.Root)
	router.GET("/sheets/default", public.Default)

	// OAuth-driven routes
	router.GET("/oauth/google", oauth.Login)
	router.GET("/oauth/google/callback", oauth.Callback)

	// Authenticated Google Sheets driven routes
	router.GET("/sheets", sheets.Index)
	router.POST("/sheets", sheets.Create)
	router.GET("/sheets/:id", sheets.Show)
	router.DELETE("/sheets/:id", sheets.Delete)

	return router
}

func main() {
	router := initRouter()
	err := router.Run(":8000")

	if err != nil {
		log.Fatal(err)
	}
}
