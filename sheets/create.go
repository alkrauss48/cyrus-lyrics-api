package sheets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alkrauss48/cyrus-lyrics-api/helpers"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func Create(c *gin.Context) {
	ctx := context.Background()
	client, err := helpers.GetGoogleOAuthClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to create sheets service")
		return
	}

	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, "Missing 'title' param")
		return
	}

	categoryLabel := "Genre"
	subCategoryLabel := "Artist"
	nameLabel := "Song Title"
	urlLabel := "URL"
	lyricsLabel := "Lyrics"
	spotifyLabel := "Spotify URL"

	resp, err := sheetsSrv.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
		Sheets: []*sheets.Sheet{
			{
				Data: []*sheets.GridData{
					{
						RowData: []*sheets.RowData{
							{
								Values: []*sheets.CellData{
									{
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: &categoryLabel,
										},
									}, {
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: &subCategoryLabel,
										},
									}, {
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: &nameLabel,
										},
									}, {
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: &urlLabel,
										},
									}, {
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: &lyricsLabel,
										},
									}, {
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: &spotifyLabel,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}).Fields("spreadsheetId").Do()

	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to create sheet: %v", err))
		return
	}

	json, err := resp.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to parse sheet json")
		return
	}

	c.Data(http.StatusCreated, gin.MIMEJSON, json)
}
