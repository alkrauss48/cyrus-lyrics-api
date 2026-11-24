package sheets

import (
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
)

// Delete moves one of the user's sheets to the Google Drive trash.
func Delete(c *gin.Context) {
	srv, err := helpers.GetDriveService(c)
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	_, err = srv.Files.Update(c.Param("id"), &drive.File{Trashed: true}).
		Context(c.Request.Context()).
		Do()
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
