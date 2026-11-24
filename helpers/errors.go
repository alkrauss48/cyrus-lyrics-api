package helpers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
)

// AbortWithError records err for the request log and responds with the
// status that StatusForError picks for it. The body names only the status, so
// upstream error details aren't exposed to clients.
//
// The iOS app signs the user out on most 4xx responses (see AGENTS.md), so
// 4xx statuses are reserved for bad tokens or bad requests.
func AbortWithError(c *gin.Context, err error) {
	status := StatusForError(err)
	_ = c.Error(err)
	c.AbortWithStatusJSON(status, gin.H{"error": http.StatusText(status)})
}

// AbortWithMessage responds with status and a client-facing message.
func AbortWithMessage(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}

// StatusForError maps an error from token parsing, OAuth, or a Google API call
// to the HTTP status to return.
func StatusForError(err error) int {
	var (
		retrieveErr *oauth2.RetrieveError
		googleErr   *googleapi.Error
		urlErr      *url.Error
	)

	switch {
	case errors.Is(err, ErrMissingToken), errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrMissingConfig):
		return http.StatusInternalServerError
	case errors.As(err, &retrieveErr):
		// Refreshing the token failed. A 4xx from Google (such as
		// invalid_grant) means the user revoked access or the refresh token
		// expired, so they need to sign in again.
		if retrieveErr.Response != nil && retrieveErr.Response.StatusCode < http.StatusInternalServerError {
			return http.StatusUnauthorized
		}
		return http.StatusBadGateway
	case errors.As(err, &googleErr):
		return statusForGoogleError(googleErr)
	case errors.As(err, &urlErr):
		// Couldn't reach Google.
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func statusForGoogleError(err *googleapi.Error) int {
	for _, item := range err.Errors {
		// Drive reports rate limiting as a 403. It's temporary, so don't
		// return a status that signs the user out.
		if item.Reason == "rateLimitExceeded" || item.Reason == "userRateLimitExceeded" {
			return http.StatusServiceUnavailable
		}
	}

	switch {
	case err.Code == http.StatusTooManyRequests:
		return http.StatusServiceUnavailable
	case err.Code >= http.StatusBadRequest && err.Code < http.StatusInternalServerError:
		return err.Code
	default:
		return http.StatusBadGateway
	}
}
