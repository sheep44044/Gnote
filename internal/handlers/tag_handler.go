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

type NoteTag struct {
	db *gorm.DB
}

func NewNoteTag(db *gorm.DB) *NoteTag {
	return &NoteTag{db: db}
}

func (h *NoteTag) GetTags(c *gin.Context) {
	var tags []models.Tag
	h.db.Find(&tags)
	utils.Success(c, tags)
}

func (h *NoteTag) GetTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := h.db.Where("id = ?", id).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tag not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "database error")
		}
		return
	}
	utils.Success(c, tag)
}

func (h *NoteTag) CreateTag(c *gin.Context) {
	var req validators.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid tag")
		return
	}

	tag := models.Tag{
		Name:  req.Name,
		Color: req.Color,
	}
	h.db.Create(&tag)
	utils.Success(c, tag)
}

func (h *NoteTag) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var req validators.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	var tag models.Tag
	if err := h.db.Where("id = ?", id).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tag not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "database error")
		}
		return
	}
	h.db.Model(&tag).Updates(models.Tag{
		Name:  req.Name,
		Color: req.Color,
	})
	utils.Success(c, tag)
}

func (h *NoteTag) DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	result := h.db.Delete(&models.Tag{}, id)
	if result.RowsAffected == 0 {
		utils.Error(c, http.StatusNotFound, "tag not found")
		return
	}
	utils.Success(c, gin.H{"message": "deleted"})
}
