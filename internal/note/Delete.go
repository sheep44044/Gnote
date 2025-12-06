package note

import (
	"log/slog"
	"net/http"
	"note/internal/cache"
	"note/internal/models"
	"note/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	cacheKeyNote := "note:" + c.Param("id")
	cacheKeyAllNotes := "notes:all"

	cache.Del(cacheKeyNote)
	cache.Del(cacheKeyAllNotes)

	slog.Info("Cache cleared for deleted note", "note_id", id)
	utils.Success(c, gin.H{"message": "deleted"})
}
