package main

import (
	"context"
	"fmt"
	"net/http"
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

	// Build the deep link to send data into the iOS app
	// baseLink = "/"
	baseLink := "cyruslyrics://login"
	deepLink := fmt.Sprintf("%s?%s", baseLink, params)

	// Redirect to the path with the parsed query params.
	c.Redirect(http.StatusFound, deepLink)
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
