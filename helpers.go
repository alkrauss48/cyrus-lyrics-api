package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func getGoogleOAuthConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	return config, nil
}

func getTokenFromRequest(c *gin.Context) (*oauth2.Token, error) {
	t, err := time.Parse(time.RFC3339, c.Query("expiry"))
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  c.Query("access_token"),
		RefreshToken: c.Query("refresh_token"),
		TokenType:    "Bearer",
		Expiry:       t,
	}, nil
}
