package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateNote_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewNoteHandler(nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateNote_EmptyTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewNoteHandler(nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"","content":"test","tags":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Title is required") {
		t.Errorf("expected 'Title is required', got %s", w.Body.String())
	}
}
