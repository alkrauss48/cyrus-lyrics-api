package oauth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// appLoginURL is the iOS app's deep link for finishing sign-in.
const appLoginURL = "cyruslyrics://login"

// Callback exchanges Google's authorization code for a token and redirects to
// the iOS app with it.
//
// The app reads only the token param: standard base64 of the JSON-encoded
// oauth2.Token. It sends that value back as a Bearer token. The
// access_token, refresh_token, and expiry params serve app versions before
// 2.0, which send them back as query params.
func Callback(c *gin.Context) {
	config, err := helpers.GetGoogleOAuthConfig()
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	authCode := c.Query("code")
	if authCode == "" {
		helpers.AbortWithMessage(c, http.StatusBadRequest, "missing 'code' param")
		return
	}

	tok, err := config.Exchange(c.Request.Context(), authCode)
	if err != nil {
		_ = c.Error(err)
		helpers.AbortWithMessage(c, http.StatusBadRequest, "unable to exchange authorization code")
		return
	}

	// Keep only the fields the API needs. Truncating the expiry to whole
	// seconds keeps the encoded token the same as earlier versions produced.
	expiry := tok.Expiry.Truncate(time.Second)
	jsonToken, err := json.Marshal(oauth2.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       expiry,
	})
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	params := url.Values{
		"token":         {base64.StdEncoding.EncodeToString(jsonToken)},
		"access_token":  {tok.AccessToken},
		"refresh_token": {tok.RefreshToken},
		"expiry":        {expiry.Format(time.RFC3339)},
	}

	c.Redirect(http.StatusFound, appLoginURL+"?"+params.Encode())
}
