// Package helpers holds the Google OAuth, token, and error-handling code
// shared by the route handlers.
package helpers

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// ErrMissingConfig means a required OAuth environment variable isn't set.
var ErrMissingConfig = errors.New("google oauth config is incomplete: CLIENT_ID, CLIENT_SECRET, and REDIRECT_URI are required")

// GetGoogleOAuthConfig builds the Google OAuth config from environment
// variables. AUTH_URI and TOKEN_URI default to Google's endpoints.
func GetGoogleOAuthConfig() (*oauth2.Config, error) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URI"),
		Scopes:       []string{drive.DriveFileScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  envOrDefault("AUTH_URI", google.Endpoint.AuthURL),
			TokenURL: envOrDefault("TOKEN_URI", google.Endpoint.TokenURL),
		},
	}

	if config.ClientID == "" || config.ClientSecret == "" || config.RedirectURL == "" {
		return nil, ErrMissingConfig
	}

	return config, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
