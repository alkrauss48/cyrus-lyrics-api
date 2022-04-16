package sheets

import (
	"context"
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func Index(c *gin.Context) {
	ctx := context.Background()

	client, err := helpers.GetGoogleOAuthClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := srv.Files.List().Fields("files(id, name)").Q("trashed = false").Do()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	json, err := resp.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to parse sheet json")
		return
	}

	c.Data(http.StatusOK, gin.MIMEJSON, json)
}
