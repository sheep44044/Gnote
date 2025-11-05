package handlers

import (
	"errors"
	"net/http"
	"note/internal/models"
	"note/internal/utils"
	"note/internal/validators"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoteHandler struct {
	db *gorm.DB
}

func NewNoteHandler(db *gorm.DB) *NoteHandler {
	return &NoteHandler{db: db}
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	var notes []models.Note
	h.db.Preload("Tags").Find(&notes) // 添加预加载
	utils.Success(c, notes)
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := h.db.Preload("Tags").Where("id = ?", id).First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "note not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "database error")
		}
		return
	}
	utils.Success(c, note)
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req validators.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid note")
		return
	}

	var tags []models.Tag
	if len(req.TagIDs) > 0 {
		h.db.Where("id IN ?", req.TagIDs).Find(&tags)
	}

	note := models.Note{
		Title:   req.Title,
		Content: req.Content,
		Tags:    tags,
	}

	h.db.Create(&note)
	utils.Success(c, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")

	var req validators.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	var note models.Note
	if err := h.db.First(&note, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "note not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "database error")
		}
		return
	}

	h.db.Model(&note).Updates(models.Note{
		Title:   req.Title,
		Content: req.Content,
	})

	var tags []models.Tag
	if len(req.TagIDs) > 0 {
		h.db.Where("id IN ?", req.TagIDs).Find(&tags)
	}
	h.db.Model(&note).Association("Tags").Replace(tags)

	h.db.Preload("Tags").First(&note, note.ID)
	utils.Success(c, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	result := h.db.Delete(&models.Note{}, id)
	if result.RowsAffected == 0 {
		utils.Error(c, http.StatusNotFound, "note not found")
		return
	}

	utils.Success(c, gin.H{"message": "deleted"})
}
