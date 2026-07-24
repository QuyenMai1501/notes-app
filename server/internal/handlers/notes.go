package handlers

import (
	"net/http"
	"strings"

	"notes-app/server/internal/database"
	"notes-app/server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteHandler struct {
	DB *pgxpool.Pool
}

func NewNoteHandler(db *pgxpool.Pool) *NoteHandler {
	return &NoteHandler{DB: db}
}

func (h *NoteHandler) List(c *gin.Context) {
	search := strings.TrimSpace(c.Query("search"))
	tag := strings.TrimSpace(c.Query("tag"))

	notes, err := database.ListNotes(c.Request.Context(), h.DB, search, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}

	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) Get(c *gin.Context) {
	id := c.Param("id")

	note, err := database.GetNote(c.Request.Context(), h.DB, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Create(c *gin.Context) {
	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if input.Tags == nil {
		input.Tags = []string{}
	}

	note, err := database.CreateNote(c.Request.Context(), h.DB, input.Title, input.Content, input.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if input.Tags == nil {
		input.Tags = []string{}
	}

	note, err := database.UpdateNote(c.Request.Context(), h.DB, id, input.Title, input.Content, input.Tags)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := database.DeleteNote(c.Request.Context(), h.DB, id)
	if err != nil {
		if strings.Contains(err.Error(), "note not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
