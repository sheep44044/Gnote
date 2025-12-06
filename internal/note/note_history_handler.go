package note

import (
	"note/internal/cache"
	"note/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// GetRecentNotes 返回最近访问的笔记ID列表（最多5个，按时间倒序）
func (h *NoteHandler) GetRecentNotes(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "未认证"})
		return
	}

	key := "user:history:" + userID

	notes, err := cache.ZRevRange(key, 0, 4)
	if err != nil {
		// Redis 出错或 key 不存在都返回空列表（更友好）
		notes = []string{}
	}
	utils.Success(c, notes)
}

// recordNoteView 记录用户访问某篇笔记（内部调用，小写开头）
func (h *NoteHandler) recordNoteView(userID, noteID string) {
	key := "user:history:" + userID
	now := float64(time.Now().Unix())

	// 1. 先移除旧记录（实现去重）
	cache.ZRem(key, noteID)

	// 2. 添加新记录（以当前时间戳为分数）
	cache.ZAdd(key, redis.Z{Score: now, Member: noteID})

	// 3. 只保留最近5条（-6 表示从第0名到倒数第6名，共删掉超出的部分）
	cache.ZRemRangeByRank(key, 0, -6)

	// 4. 设置30天自动过期（可选但推荐）
	cache.Expire(key, 30*24*time.Hour)
}
