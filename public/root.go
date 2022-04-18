package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.String(http.StatusOK, "TEST 1.0.3: You've found the CyrusLyrics API. "+
		"For more information over the CyrusLyrics iOS app, navigate "+
		"to https://cyruskrauss.com")
}
