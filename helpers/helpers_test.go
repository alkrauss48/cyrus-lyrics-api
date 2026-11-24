package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
)

func newContext(t *testing.T, target string, authHeader string) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, target, nil)
	if authHeader != "" {
		c.Request.Header.Set("Authorization", authHeader)
	}
	return c
}

func encodeToken(t *testing.T, tok oauth2.Token) string {
	t.Helper()

	b, err := json.Marshal(tok)
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func TestGetTokenFromRequestHeaders(t *testing.T) {
	future := time.Now().Add(time.Hour).Truncate(time.Second)
	past := time.Now().Add(-time.Hour)

	tests := []struct {
		name    string
		header  string
		wantErr error
	}{
		{"missing header", "", ErrMissingToken},
		{"not bearer", "Basic abc", ErrInvalidToken},
		{"bearer with no token", "Bearer ", ErrInvalidToken},
		{"not base64", "Bearer !!!", ErrInvalidToken},
		{"not json", "Bearer " + base64.StdEncoding.EncodeToString([]byte("nope")), ErrInvalidToken},
		{"no access token", "Bearer " + encodeToken(t, oauth2.Token{RefreshToken: "r"}), ErrInvalidToken},
		{"expired without refresh token", "Bearer " + encodeToken(t, oauth2.Token{AccessToken: "a", Expiry: past}), ErrInvalidToken},
		{"expired with refresh token", "Bearer " + encodeToken(t, oauth2.Token{AccessToken: "a", RefreshToken: "r", Expiry: past}), nil},
		{"valid", "Bearer " + encodeToken(t, oauth2.Token{AccessToken: "a", RefreshToken: "r", Expiry: future}), nil},
		{"lowercase scheme", "bearer " + encodeToken(t, oauth2.Token{AccessToken: "a", RefreshToken: "r", Expiry: future}), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok, err := GetTokenFromRequest(newContext(t, "/", tt.header))

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && tok.AccessToken != "a" {
				t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "a")
			}
		})
	}
}

func TestGetTokenFromRequestQuery(t *testing.T) {
	expiry := time.Now().Add(time.Hour).UTC().Truncate(time.Second)

	t.Run("valid", func(t *testing.T) {
		q := url.Values{
			"access_token":  {"a"},
			"refresh_token": {"r"},
			"expiry":        {expiry.Format(time.RFC3339)},
		}

		tok, err := GetTokenFromRequest(newContext(t, "/?"+q.Encode(), ""))
		if err != nil {
			t.Fatal(err)
		}
		if tok.AccessToken != "a" || tok.RefreshToken != "r" || !tok.Expiry.Equal(expiry) {
			t.Errorf("token = %+v", tok)
		}
	})

	t.Run("bad expiry", func(t *testing.T) {
		_, err := GetTokenFromRequest(newContext(t, "/?access_token=a&expiry=tomorrow", ""))
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidToken)
		}
	})
}

func TestGetGoogleOAuthConfig(t *testing.T) {
	t.Run("missing values", func(t *testing.T) {
		t.Setenv("CLIENT_ID", "")
		t.Setenv("CLIENT_SECRET", "")
		t.Setenv("REDIRECT_URI", "")

		if _, err := GetGoogleOAuthConfig(); !errors.Is(err, ErrMissingConfig) {
			t.Fatalf("err = %v, want %v", err, ErrMissingConfig)
		}
	})

	t.Run("defaults to Google endpoints", func(t *testing.T) {
		t.Setenv("CLIENT_ID", "id")
		t.Setenv("CLIENT_SECRET", "secret")
		t.Setenv("REDIRECT_URI", "https://example.com/oauth/google/callback")
		t.Setenv("AUTH_URI", "")
		t.Setenv("TOKEN_URI", "")

		config, err := GetGoogleOAuthConfig()
		if err != nil {
			t.Fatal(err)
		}
		if config.Endpoint.AuthURL != "https://accounts.google.com/o/oauth2/auth" {
			t.Errorf("AuthURL = %q", config.Endpoint.AuthURL)
		}
		if config.Endpoint.TokenURL != "https://oauth2.googleapis.com/token" {
			t.Errorf("TokenURL = %q", config.Endpoint.TokenURL)
		}
		if len(config.Scopes) != 1 || config.Scopes[0] != "https://www.googleapis.com/auth/drive.file" {
			t.Errorf("Scopes = %v", config.Scopes)
		}
	})
}

func TestStatusForError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"missing token", ErrMissingToken, http.StatusUnauthorized},
		{"invalid token", fmt.Errorf("%w: detail", ErrInvalidToken), http.StatusUnauthorized},
		{"missing config", ErrMissingConfig, http.StatusInternalServerError},
		{
			"refresh rejected",
			&url.Error{Op: "Get", URL: "x", Err: &oauth2.RetrieveError{Response: &http.Response{StatusCode: http.StatusBadRequest}}},
			http.StatusUnauthorized,
		},
		{
			"refresh server error",
			&oauth2.RetrieveError{Response: &http.Response{StatusCode: http.StatusInternalServerError}},
			http.StatusBadGateway,
		},
		{"google unauthorized", &googleapi.Error{Code: http.StatusUnauthorized}, http.StatusUnauthorized},
		{"google not found", &googleapi.Error{Code: http.StatusNotFound}, http.StatusNotFound},
		{
			"google rate limit as 403",
			&googleapi.Error{Code: http.StatusForbidden, Errors: []googleapi.ErrorItem{{Reason: "userRateLimitExceeded"}}},
			http.StatusServiceUnavailable,
		},
		{"google 429", &googleapi.Error{Code: http.StatusTooManyRequests}, http.StatusServiceUnavailable},
		{"google 503", &googleapi.Error{Code: http.StatusServiceUnavailable}, http.StatusBadGateway},
		{"network", &url.Error{Op: "Get", URL: "x", Err: errors.New("connection refused")}, http.StatusBadGateway},
		{"unknown", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusForError(tt.err); got != tt.want {
				t.Errorf("StatusForError() = %d, want %d", got, tt.want)
			}
		})
	}
}
