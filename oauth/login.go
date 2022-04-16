package oauth

import (
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Login(c *gin.Context) {
	config, err := helpers.GetGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	authURL := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	c.Redirect(http.StatusFound, authURL)
}
