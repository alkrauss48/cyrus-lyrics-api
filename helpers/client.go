package helpers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGoogleOAuthClient(c *gin.Context) (*http.Client, error) {
	tok, err := GetTokenFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, err
	}

	config, err := GetGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, err
	}

	client := config.Client(context.Background(), tok)

	return client, nil
}
