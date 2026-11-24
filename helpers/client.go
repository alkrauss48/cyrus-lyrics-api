package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// GetGoogleOAuthClient returns an HTTP client authorized with the request's
// token. It refreshes the token when needed.
func GetGoogleOAuthClient(c *gin.Context) (*http.Client, error) {
	tok, err := GetTokenFromRequest(c)
	if err != nil {
		return nil, err
	}

	config, err := GetGoogleOAuthConfig()
	if err != nil {
		return nil, err
	}

	return config.Client(c.Request.Context(), tok), nil
}

// GetDriveService returns a Google Drive client for the request's user.
func GetDriveService(c *gin.Context) (*drive.Service, error) {
	client, err := GetGoogleOAuthClient(c)
	if err != nil {
		return nil, err
	}

	return drive.NewService(c.Request.Context(), option.WithHTTPClient(client))
}

// GetSheetsService returns a Google Sheets client for the request's user.
func GetSheetsService(c *gin.Context) (*sheets.Service, error) {
	client, err := GetGoogleOAuthClient(c)
	if err != nil {
		return nil, err
	}

	return sheets.NewService(c.Request.Context(), option.WithHTTPClient(client))
}
