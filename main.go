package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
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

	params := fmt.Sprintf("access_token=%s&refresh_token=%s", tok.AccessToken, tok.RefreshToken)

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

// func listDriveFiles(c *gin.Context) {
//	ctx := context.Background()
// 	// NOTE: Get Token
// 	client := config.Client(context.Background(), tok)
//
// 	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Drive client: %v", err)
// 	}
//
// 	r, err := srv.Files.List().PageSize(10).
// 		Fields("nextPageToken, files(id, name)").Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve files: %v", err)
// 	}
// 	fmt.Println("Files:")
// 	if len(r.Files) == 0 {
// 		fmt.Println("No files found.")
// 	} else {
// 		for _, i := range r.Files {
// 			fmt.Printf("%s (%s)\n", i.Name, i.Id)
// 		}
// 	}
// }

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
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
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

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
