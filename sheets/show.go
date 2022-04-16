package sheets

import (
	"context"
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func Show(c *gin.Context) {
	ctx := context.Background()
	spreadsheetId := c.Param("id")

	client, err := helpers.GetGoogleOAuthClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Get all A-F cells in the first visible sheet
	// Google Sheets A1 Range Notation
	sheetRange := "A:F"

	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to create sheets service")
		return
	}

	resp, err := sheetsSrv.Spreadsheets.Values.Get(spreadsheetId, sheetRange).Context(ctx).Do()
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to read sheet")
		return
	}

	json, err := resp.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to parse sheet json")
		return
	}

	c.Data(http.StatusOK, gin.MIMEJSON, json)
}
