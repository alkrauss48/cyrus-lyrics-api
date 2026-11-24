package sheets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/sheets", Index)
	router.POST("/sheets", Create)
	router.GET("/sheets/:id", Show)
	router.DELETE("/sheets/:id", Delete)
	return router
}

// The iOS app signs the user out on a 401, so a missing or malformed token
// must return 401 before any call to Google.
func TestRoutesRequireToken(t *testing.T) {
	router := newRouter()

	routes := []struct{ method, target string }{
		{http.MethodGet, "/sheets"},
		{http.MethodPost, "/sheets?title=My+Songs"},
		{http.MethodGet, "/sheets/abc"},
		{http.MethodDelete, "/sheets/abc"},
	}

	for _, header := range []string{"", "Bearer not-a-token"} {
		for _, route := range routes {
			req := httptest.NewRequest(route.method, route.target, nil)
			if header != "" {
				req.Header.Set("Authorization", header)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("%s %s with header %q: status = %d, want %d",
					route.method, route.target, header, w.Code, http.StatusUnauthorized)
			}
		}
	}
}

func TestCreateRequiresTitle(t *testing.T) {
	w := httptest.NewRecorder()
	newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/sheets", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestNewSpreadsheetHeaderRow(t *testing.T) {
	s := newSpreadsheet("My Songs")

	if s.Properties.Title != "My Songs" {
		t.Errorf("title = %q", s.Properties.Title)
	}

	cells := s.Sheets[0].Data[0].RowData[0].Values
	want := []string{"Genre", "Artist", "Song Title", "URL", "Lyrics", "Spotify URL"}
	if len(cells) != len(want) {
		t.Fatalf("got %d header cells, want %d", len(cells), len(want))
	}
	for i, label := range want {
		if got := *cells[i].UserEnteredValue.StringValue; got != label {
			t.Errorf("header %d = %q, want %q", i, got, label)
		}
	}
}
