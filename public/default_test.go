package public

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestDefaultRoute(t *testing.T) {
	router := gin.Default()
	router.GET("/", Default)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Convert the JSON response to a map
	var response []DefaultFile
	_ = json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Equal(t, len(response), 2)
}
