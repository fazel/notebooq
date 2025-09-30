package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/fazel/notebooq/internal/models"
	"github.com/fazel/notebooq/internal/repository"
	"github.com/fazel/notebooq/internal/service"
)

type NoteHandler struct {
	svc *service.NoteService
}

// NewNoteHandler constructor
func NewNoteHandler(db *gorm.DB) *NoteHandler {
	repo := repository.NewNoteRepo(db)
	svc := service.NewNoteService(repo)
	return &NoteHandler{svc: svc}
}

// CreateNote → POST /notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := &models.Note{
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.svc.Create(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// ListNotes → GET /notes
func (h *NoteHandler) ListNotes(c *gin.Context) {
	userID := c.GetUint("userID")

	notes, err := h.svc.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// GetNote → GET /notes/:id
func (h *NoteHandler) GetNote(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	note, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// UpdateNote → PUT /notes/:id
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	note.Title = req.Title
	note.Content = req.Content
	note.UpdatedAt = time.Now()

	if err := h.svc.Update(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote → DELETE /notes/:id
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}
