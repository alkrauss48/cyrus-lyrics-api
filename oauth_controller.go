package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func googleLoginCallback(c *gin.Context) {
	config, err := getGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	authCode := c.Query("code")
	if authCode == "" {
		c.JSON(http.StatusBadRequest, "Missing 'code' param")
	}

	// Parse the OAuth authorization token for its individual parts,
	// and then send it back to a new route as query params so that it can be
	// accessed from the browser.
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to retrieve token from web")
		return
	}

	// Build the query parameters for the parsed token
	params := fmt.Sprintf(
		"access_token=%s&refresh_token=%s&expiry=%s",
		tok.AccessToken,
		tok.RefreshToken,
		tok.Expiry.Format(time.RFC3339),
	)

	// Build the full path for the parsed token
	location := url.URL{
		Path:     "/oauth/google/processed",
		RawQuery: params,
	}

	// Redirect to the path with the parsed query params.
	c.Redirect(http.StatusFound, location.RequestURI())
}

// Don't need to show anything here; this route is just for the mobile app to
// get the parsed token information and store it.
func googleLoginProcessed(c *gin.Context) {
	c.JSON(http.StatusOK, "success")
}

func googleLogin(c *gin.Context) {
	config, err := getGoogleOAuthConfig()
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
