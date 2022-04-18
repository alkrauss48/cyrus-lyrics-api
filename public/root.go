package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.String(http.StatusOK, "Hello, OKC WebDevs! "+
		"To see our talk slides, visit https://clvr.sh/okcwebdevs.")
}
