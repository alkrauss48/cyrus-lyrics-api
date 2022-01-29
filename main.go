package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", root)
	router.GET("/oauth/google", googleLogin)
	router.GET("/oauth/google/callback", googleLoginCallback)
	router.GET("/oauth/google/processed", googleLoginProcessed)
	router.GET("/sheets", getAllSheets)
	router.GET("/sheets/new", newSheet)
	router.GET("/sheets/default", getDefaultSheetIds)
	router.GET("/sheets/:id", getSheetByID)

	router.Run(":8000")
}

func root(c *gin.Context) {
	c.String(http.StatusOK, "You've found the CyrusLyrics API. "+
		"For more information over the CyrusLyrics iOS app, navigate "+
		"to https://cyruskrauss.com")
}
