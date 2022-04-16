package helpers

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GetTokenFromRequest(c *gin.Context) (*oauth2.Token, error) {
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
