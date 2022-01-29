package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GoogleConfig struct {
	Web SubGoogleConfig `json:"web"`
}

type SubGoogleConfig struct {
	ClientId                string    `json:"client_id"`
	ProjectId               string    `json:"project_id"`
	AuthUri                 string    `json:"auth_uri"`
	TokenUri                string    `json:"token_uri"`
	AuthProviderX509CertUrl string    `json:"auth_provider_x509_cert_url"`
	ClientSecret            string    `json:"client_secret"`
	RedirectUris            [1]string `json:"redirect_uris"`
}

func marshalOAuthCredentials() ([]byte, error) {
	credentials := GoogleConfig{
		Web: SubGoogleConfig{
			AuthProviderX509CertUrl: os.Getenv("AUTH_PROVIDER_CERT_URL"),
			AuthUri:                 os.Getenv("AUTH_URI"),
			ClientId:                os.Getenv("CLIENT_ID"),
			ClientSecret:            os.Getenv("CLIENT_SECRET"),
			ProjectId:               os.Getenv("PROJECT_ID"),
			RedirectUris:            [1]string{os.Getenv("REDIRECT_URI")},
			TokenUri:                os.Getenv("TOKEN_URI"),
		},
	}

	jsonCredentials, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	return jsonCredentials, nil
}

func getGoogleOAuthConfig() (*oauth2.Config, error) {
	jsonCredentials, err := marshalOAuthCredentials()
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(jsonCredentials, drive.DriveFileScope)
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

func getGoogleOAuthClient(c *gin.Context) (*http.Client, error) {
	tok, err := getTokenFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, err
	}

	config, err := getGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, err
	}

	client := config.Client(context.Background(), tok)

	return client, nil
}
