package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type DefaultFile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func getDefaultSheetIds(c *gin.Context) {
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

func getAllSheets(c *gin.Context) {
	ctx := context.Background()

	client, err := getGoogleOAuthClient(c)
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

func getSheetByID(c *gin.Context) {
	ctx := context.Background()
	spreadsheetId := c.Param("id")

	client, err := getGoogleOAuthClient(c)
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

func newSheet(c *gin.Context) {
	ctx := context.Background()
	client, err := getGoogleOAuthClient(c)
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

	categoryLabel := "Category"
	subCategoryLabel := "Subcategory"
	nameLabel := "Name"
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
