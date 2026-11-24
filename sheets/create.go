package sheets

import (
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
)

// headerLabels is the header row of a new sheet. The iOS app reads the
// columns by position: category, sub-category, song name, URL, lyrics,
// and Spotify URL.
var headerLabels = []string{"Genre", "Artist", "Song Title", "URL", "Lyrics", "Spotify URL"}

// Create makes a new sheet named by the title query param, with the header
// row filled in, and returns its spreadsheetId.
func Create(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		helpers.AbortWithMessage(c, http.StatusBadRequest, "missing 'title' param")
		return
	}

	srv, err := helpers.GetSheetsService(c)
	if err != nil {
		helpers.AbortWithError(c, err)
		return
	}

	resp, err := srv.Spreadsheets.Create(newSpreadsheet(title)).
		Fields("spreadsheetId").
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

	c.Data(http.StatusCreated, gin.MIMEJSON, body)
}

func newSpreadsheet(title string) *sheets.Spreadsheet {
	cells := make([]*sheets.CellData, len(headerLabels))
	for i, label := range headerLabels {
		cells[i] = &sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{StringValue: &label},
		}
	}

	return &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{Title: title},
		Sheets: []*sheets.Sheet{{
			Data: []*sheets.GridData{{
				RowData: []*sheets.RowData{{Values: cells}},
			}},
		}},
	}
}
