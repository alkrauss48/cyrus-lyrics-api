package public

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func serve(t *testing.T, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/", handler)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	return w
}

func TestRoot(t *testing.T) {
	w := serve(t, Root)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDefault(t *testing.T) {
	t.Setenv("DEMO_SHEET_ID", "demo-id")
	t.Setenv("AARON_SHEET_ID", "aaron-id")

	w := serve(t, Default)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got []DefaultFile
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	want := []DefaultFile{
		{ID: "demo-id", Name: "Demo Songs"},
		{ID: "aaron-id", Name: "Cyrus' Dad's Songs"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d files, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("file %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}
