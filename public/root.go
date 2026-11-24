// Package public has the unauthenticated routes.
package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Root responds with a plain-text placeholder message.
func Root(c *gin.Context) {
	c.String(http.StatusOK, "Hello Innotech! To register for this talk's course, "+
		"navigate to https://clvr.sh/innotech2022")
}
