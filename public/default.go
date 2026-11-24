package public

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// DefaultFile is a preloaded, publicly readable sheet.
type DefaultFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Default lists the preloaded sheets that every user can see. Their IDs come
// from the DEMO_SHEET_ID and AARON_SHEET_ID environment variables.
func Default(c *gin.Context) {
	c.JSON(http.StatusOK, []DefaultFile{
		{
			ID:   os.Getenv("DEMO_SHEET_ID"),
			Name: "Demo Songs",
		}, {
			ID:   os.Getenv("AARON_SHEET_ID"),
			Name: "Cyrus' Dad's Songs",
		},
	})
}
