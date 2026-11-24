// Package oauth handles Google sign-in for the iOS app.
package oauth

import (
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// Login redirects to Google's consent screen. It requests offline access and
// forces the consent prompt, so Google always returns a refresh token.
func Login(c *gin.Context) {
	config, err := helpers.GetGoogleOAuthConfig()
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	authURL := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	c.Redirect(http.StatusFound, authURL)
}
