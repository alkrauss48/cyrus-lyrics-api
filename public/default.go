package public

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type DefaultFile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func Default(c *gin.Context) {
	c.JSON(http.StatusOK, []DefaultFile{
		{
			Id:   os.Getenv("DEMO_SHEET_ID"),
			Name: "Demo Songs",
		}, {
			Id:   os.Getenv("AARON_SHEET_ID"),
			Name: "Cyrus' Dad's Songs",
		},
	})
}
