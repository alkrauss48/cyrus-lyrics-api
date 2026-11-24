package sheets

import (
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
)

// sheetRange is every row of columns A-F in the first visible sheet, in
// Google Sheets A1 notation.
const sheetRange = "A:F"

// Show returns the values in one of the user's sheets.
func Show(c *gin.Context) {
	srv, err := helpers.GetSheetsService(c)
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	resp, err := srv.Spreadsheets.Values.
		Get(c.Param("id"), sheetRange).
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
