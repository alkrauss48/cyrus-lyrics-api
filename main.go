package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getGoogleOAuthConfig() *oauth2.Config {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// config, err := google.ConfigFromJSON(b, drive.DriveScope)
	// config, err := google.ConfigFromJSON(b, drive.DriveFileScope, drive.DriveReadonlyScope)
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config
}

func googleLoginCallback(c *gin.Context) {
	config := getGoogleOAuthConfig()
	authCode := c.Query("code")

	// Parse the OAuth authorization token for its individual parts,
	// and then send it back to a new route as query params so that it can be
	// accessed from the browser.
	tok, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	params := fmt.Sprintf(
		"access_token=%s&refresh_token=%s&expiry=%s",
		tok.AccessToken,
		tok.RefreshToken,
		tok.Expiry,
	)

	location := url.URL{
		Path:     "/oauth/google/processed",
		RawQuery: params,
	}

	c.Redirect(http.StatusFound, location.RequestURI())
}

func googleLoginProcessed(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "success")
}

func googleLogin(c *gin.Context) {
	config := getGoogleOAuthConfig()

	authURL := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	c.Redirect(http.StatusFound, authURL)
}

func getTokenFromRequest(c *gin.Context) *oauth2.Token {
	t, err := time.Parse(c.Query("expiry"), c.Query("expiry"))
	if err != nil {
		// TODO: Handle error.
	}

	return &oauth2.Token{
		AccessToken:  c.Query("access_token"),
		RefreshToken: c.Query("refresh_token"),
		TokenType:    "Bearer",
		Expiry:       t,
	}
}

func copyEmptySheet(c *gin.Context) {
	ctx := context.Background()
	tok := getTokenFromRequest(c)
	config := getGoogleOAuthConfig()

	client := config.Client(context.Background(), tok)

	driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	categoryLabel := "Category"
	subCategoryLabel := "Subcategory"
	nameLabel := "Name"
	urlLabel := "URL"
	lyricsLabel := "Lyrics"
	spotifyLabel := "Spotify URL"

	_, err = sheetsSrv.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "New Sheet From API",
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
		log.Fatalf("Unable to create sheet: %v", err)
	}

	r, err := driveSrv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default()
	router.GET("/oauth/google", googleLogin)
	router.GET("/oauth/google/callback", googleLoginCallback)
	router.GET("/oauth/google/processed", googleLoginProcessed)
	router.GET("/sheets/copy/empty", copyEmptySheet)

	// TODO: Remove these
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8000")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
