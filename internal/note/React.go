package note

import (
	"net/http"
	"note/internal/models"
	"note/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *NoteHandler) ReactToNote(c *gin.Context) {
	noteID := c.Param("id")
	noteIDUint64, _ := strconv.ParseUint(noteID, 10, 64)
	userid, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	userID, ok := userid.(uint)
	if !ok {
		utils.Error(c, http.StatusInternalServerError, "用户ID类型错误")
		return
	}

	var input struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "需要 emoji")
		return
	}

	// 校验 emoji（简单白名单）
	validEmojis := map[string]bool{
		"❤️": true, "👍": true, "🔥": true, "👏": true, "😂": true, "😮": true,
	}
	if !validEmojis[input.Emoji] {
		utils.Error(c, http.StatusBadRequest, "不支持的 emoji")
		return
	}

	// 删除用户对该笔记的旧 reaction（同一时间只能点一个）
	h.db.Where("user_id = ? AND note_id = ?", userID, noteID).Delete(&models.Reaction{})

	reaction := models.Reaction{
		UserID: userID,
		NoteID: uint(noteIDUint64),
		Emoji:  input.Emoji,
	}
	h.db.Create(&reaction)

	utils.Success(c, gin.H{"message": "反应成功"})
}
