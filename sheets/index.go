// Package sheets has the authenticated routes for the user's own Google Sheets.
package sheets

import (
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
)

// Index lists the IDs and names of the user's sheets that aren't trashed.
// Because of the drive.file scope, these are only sheets this app created.
func Index(c *gin.Context) {
	srv, err := helpers.GetDriveService(c)
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	resp, err := srv.Files.List().
		Fields("files(id, name)").
		Q("trashed = false").
		Context(c.Request.Context()).
		Do()
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	body, err := resp.MarshalJSON()
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	c.Data(http.StatusOK, gin.MIMEJSON, body)
}
