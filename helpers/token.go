package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	// ErrMissingToken means the request didn't include any credentials.
	ErrMissingToken = errors.New("missing bearer token")

	// ErrInvalidToken means the request's credentials couldn't be used.
	ErrInvalidToken = errors.New("invalid bearer token")
)

const bearerPrefix = "Bearer "

// HasTokenQuery reports whether the request uses the legacy query-param auth.
func HasTokenQuery(c *gin.Context) bool {
	return c.Query("access_token") != ""
}

// GetTokenFromRequestQuery reads a token from the access_token,
// refresh_token, and expiry query params.
//
// Deprecated: only app versions before 2.0 send credentials this way. Use
// [GetTokenFromRequestHeaders].
func GetTokenFromRequestQuery(c *gin.Context) (*oauth2.Token, error) {
	expiry, err := time.Parse(time.RFC3339, c.Query("expiry"))
	if err != nil {
		return nil, fmt.Errorf("%w: unparseable expiry: %w", ErrInvalidToken, err)
	}

	return validateToken(&oauth2.Token{
		AccessToken:  c.Query("access_token"),
		RefreshToken: c.Query("refresh_token"),
		TokenType:    "Bearer",
		Expiry:       expiry,
	})
}

// GetTokenFromRequestHeaders reads a token from an
// "Authorization: Bearer <token>" header, where <token> is standard base64
// of a JSON-encoded [oauth2.Token], as created by the OAuth callback.
func GetTokenFromRequestHeaders(c *gin.Context) (*oauth2.Token, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, ErrMissingToken
	}

	if len(authHeader) <= len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
		return nil, fmt.Errorf("%w: expected a Bearer authorization header", ErrInvalidToken)
	}

	decodedToken, err := base64.StdEncoding.DecodeString(authHeader[len(bearerPrefix):])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(decodedToken, &token); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	return validateToken(&token)
}

// GetTokenFromRequest reads the request's token, preferring the legacy query
// params when they're present.
func GetTokenFromRequest(c *gin.Context) (*oauth2.Token, error) {
	if HasTokenQuery(c) {
		return GetTokenFromRequestQuery(c)
	}

	return GetTokenFromRequestHeaders(c)
}

// validateToken rejects tokens that can never authenticate: ones with no
// access token, and expired ones with no refresh token. Without this check,
// those failures would surface later as opaque transport errors.
func validateToken(token *oauth2.Token) (*oauth2.Token, error) {
	if token.AccessToken == "" {
		return nil, fmt.Errorf("%w: no access token", ErrInvalidToken)
	}

	if token.RefreshToken == "" && !token.Valid() {
		return nil, fmt.Errorf("%w: expired with no refresh token", ErrInvalidToken)
	}

	return token, nil
}
