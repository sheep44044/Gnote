package note

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"note/internal/models"
	"note/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// GetRecentNotes 返回最近访问的笔记ID列表（最多5个，按时间倒序）
func (h *NoteHandler) GetRecentNotes(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	key := fmt.Sprintf("user:history:%d", userID)

	noteIDs, err := h.cache.ZRevRange(c, key, 0, 4)
	if err != nil {
		// Redis 出错或 key 不存在都返回空列表（更友好）
		noteIDs = []string{}
	}

	if len(noteIDs) == 0 {
		var histories []models.History
		if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(5).Find(&histories).Error; err == nil {
			// 修正点 1：使用 append 避免 Panic
			for _, h := range histories {
				noteIDs = append(noteIDs, strconv.Itoa(int(h.NoteID)))
			}
		}

		if len(noteIDs) == 0 {
			utils.Success(c, []interface{}{})
			return
		}

		var notes []models.Note
		if err := h.db.Where("id = ?", noteIDs).Find(&notes).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, err.Error())
		}

		noteMap := make(map[uint]models.Note)
		for _, n := range notes {
			noteMap[n.ID] = n
		}

		type NoteDTO struct {
			ID            uint      `json:"id"`
			Title         string    `json:"title"`
			Content       string    `json:"content"`
			FavoriteCount int       `json:"favorite_count"`
			IsFavorite    bool      `json:"is_favorite"`
			CreatedAt     time.Time `json:"created_at"`
		}
		result := make([]NoteDTO, 0, len(noteIDs))

		for _, idStr := range noteIDs {
			idUint64, _ := strconv.ParseUint(idStr, 10, 32)
			idUint := uint(idUint64)

			if note, exists := noteMap[idUint]; exists {
				result = append(result, NoteDTO{
					ID:            note.ID,
					Title:         note.Title,
					Content:       note.Content,
					FavoriteCount: note.FavoriteCount,
					IsFavorite:    note.IsFavorite,
					CreatedAt:     note.CreatedAt,
				})
			}
		}
		utils.Success(c, result)
	}
}

// 这里的入参类型是 uint，确保调用方传递正确
func (h *NoteHandler) recordNoteView(ctx context.Context, userID, noteID uint) {
	key := fmt.Sprintf("user:history:%d", userID)
	now := float64(time.Now().Unix())

	// Redis ZSet 的 Member 建议存 String，虽然有些库支持 Int，但明确转 String 更安全
	noteIDStr := strconv.Itoa(int(noteID))

	// 1. 先移除旧记录
	h.cache.ZRem(ctx, key, noteIDStr)

	// 2. 添加新记录
	h.cache.ZAdd(ctx, key, redis.Z{Score: now, Member: noteIDStr})

	// 3. 截断
	h.cache.ZRemRangeByRank(ctx, key, 0, -6)

	// 4. 过期时间
	h.cache.Expire(ctx, key, 30*24*time.Hour)

	// 5. 异步发送
	msg := models.HistoryMsg{UserID: userID, NoteID: noteID}
	body, _ := json.Marshal(msg)
	// 确保 h.rabbit 不为 nil
	if h.rabbit != nil {
		h.rabbit.Publish("history_queue", body)
	}
}
