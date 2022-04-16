package sheets

import (
	"context"
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func Delete(c *gin.Context) {
	ctx := context.Background()
	spreadsheetId := c.Param("id")

	client, err := helpers.GetGoogleOAuthClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = srv.Files.Update(spreadsheetId, &drive.File{
		Trashed: true,
	}).Do()

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Status(http.StatusNoContent)
}
