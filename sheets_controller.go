package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getSheetByID(c *gin.Context) {
	spreadsheetId := c.Param("id")

	ctx := context.Background()
	tok, err := getTokenFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	config, err := getGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	client := config.Client(context.Background(), tok)

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

	c.JSON(http.StatusOK, string(json))
}

func newSheet(c *gin.Context) {
	ctx := context.Background()
	tok, err := getTokenFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	config, err := getGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	client := config.Client(context.Background(), tok)

	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to create sheets service")
		return
	}

	categoryLabel := "Category"
	subCategoryLabel := "Subcategory"
	nameLabel := "Name"
	urlLabel := "URL"
	lyricsLabel := "Lyrics"
	spotifyLabel := "Spotify URL"

	_, err = sheetsSrv.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: c.Query("title"),
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
	}).Do()

	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to create sheet: %v", err))
		return
	}

	c.JSON(http.StatusCreated, nil)
}
