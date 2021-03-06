package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func HasTokenQuery(c *gin.Context) bool {
	return c.Query("access_token") != ""
}

// Deprecated
func GetTokenFromRequestQuery(c *gin.Context) (*oauth2.Token, error) {
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

func GetTokenFromRequestHeaders(c *gin.Context) (*oauth2.Token, error) {
	const BEARER_SCHEMA = "Bearer "

	authHeader := c.GetHeader("Authorization")
	if len(authHeader) <= len(BEARER_SCHEMA) {
		return nil, errors.New("Invalid token")
	}

	tokenString := authHeader[len(BEARER_SCHEMA):]

	decodedToken, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}

	var token = oauth2.Token{}
	err = json.Unmarshal(decodedToken, &token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GetTokenFromRequest(c *gin.Context) (*oauth2.Token, error) {
	if HasTokenQuery(c) {
		return GetTokenFromRequestQuery(c)
	}

	return GetTokenFromRequestHeaders(c)
}
