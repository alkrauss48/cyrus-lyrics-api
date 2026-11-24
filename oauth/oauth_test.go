package oauth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func setConfig(t *testing.T, tokenURL string) {
	t.Helper()
	t.Setenv("CLIENT_ID", "client-id")
	t.Setenv("CLIENT_SECRET", "client-secret")
	t.Setenv("REDIRECT_URI", "https://api.example.com/oauth/google/callback")
	t.Setenv("AUTH_URI", "")
	t.Setenv("TOKEN_URI", tokenURL)
}

func serve(t *testing.T, handler gin.HandlerFunc, target string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/", handler)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
	return w
}

func TestLoginRedirectsToGoogle(t *testing.T) {
	setConfig(t, "")

	w := serve(t, Login, "/")

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	loc, err := url.Parse(w.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if loc.Host != "accounts.google.com" {
		t.Errorf("redirect host = %q", loc.Host)
	}

	q := loc.Query()
	for key, want := range map[string]string{
		"client_id":   "client-id",
		"access_type": "offline",
		"prompt":      "consent",
		"scope":       "https://www.googleapis.com/auth/drive.file",
	} {
		if got := q.Get(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestLoginWithoutConfig(t *testing.T) {
	t.Setenv("CLIENT_ID", "")

	if w := serve(t, Login, "/"); w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestCallbackMissingCode(t *testing.T) {
	setConfig(t, "")

	if w := serve(t, Callback, "/"); w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCallbackExchangeFails(t *testing.T) {
	google := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer google.Close()
	setConfig(t, google.URL)

	if w := serve(t, Callback, "/?code=bad"); w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestCallbackRedirectsToApp checks the deep link that the iOS app parses.
// The token values include characters that must be escaped in a query string.
func TestCallbackRedirectsToApp(t *testing.T) {
	const (
		accessToken  = "ya29.a+b/c="
		refreshToken = "1//0refresh+token"
	)

	google := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil || r.Form.Get("code") != "the-code" {
			http.Error(w, "bad code", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}))
	defer google.Close()
	setConfig(t, google.URL)

	w := serve(t, Callback, "/?code=the-code&state=state-token")

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusFound, w.Body)
	}

	loc, err := url.Parse(w.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if loc.Scheme != "cyruslyrics" || loc.Host != "login" {
		t.Errorf("redirect = %q, want cyruslyrics://login", loc)
	}

	q := loc.Query()
	if q.Get("access_token") != accessToken || q.Get("refresh_token") != refreshToken {
		t.Errorf("legacy params = %v", q)
	}

	decoded, err := base64.StdEncoding.DecodeString(q.Get("token"))
	if err != nil {
		t.Fatalf("token isn't standard base64: %v", err)
	}

	var tok oauth2.Token
	if err := json.Unmarshal(decoded, &tok); err != nil {
		t.Fatalf("token isn't a JSON oauth2.Token: %v", err)
	}
	if tok.AccessToken != accessToken || tok.RefreshToken != refreshToken || tok.TokenType != "Bearer" {
		t.Errorf("token = %+v", tok)
	}
	if tok.Expiry.IsZero() || tok.Expiry.Nanosecond() != 0 {
		t.Errorf("expiry = %v, want whole seconds", tok.Expiry)
	}
	if got := q.Get("expiry"); got != tok.Expiry.Format("2006-01-02T15:04:05Z07:00") {
		t.Errorf("expiry param = %q, token expiry = %v", got, tok.Expiry)
	}
}
