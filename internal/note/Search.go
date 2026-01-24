package note

import (
	"log/slog"
	"net/http"
	"note/internal/models"
	"note/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *NoteHandler) SearchNotes(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// 1. 获取查询参数
	query := c.Query("q")
	if query == "" {
		utils.Error(c, http.StatusBadRequest, "缺少搜索关键词 'q'")
		return
	}

	// 2. 安全限制：关键词长度 ≤ 50 字符（防滥用）
	if len(query) > 50 {
		utils.Error(c, http.StatusBadRequest, "搜索词过长")
		return
	}

	// 3. 分页（可选）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 10

	// 4. 构建 LIKE 查询（GORM 自动转义，防注入）
	var notes []models.Note
	offset := (page - 1) * pageSize

	// 同时搜标题和内容（忽略大小写）
	err = h.db.Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").
		Where("user_id = ?", userID). // ← 别忘了权限！
		Order("updated_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&notes).Error

	if err != nil {
		slog.Error("Search notes failed", "error", err)
		utils.Error(c, http.StatusInternalServerError, "搜索失败")
		return
	}

	utils.Success(c, gin.H{
		"notes": notes,
		"page":  page,
		"total": len(notes),
	})
}

func (h *NoteHandler) SmartSearch(c *gin.Context) {
	query := c.Query("q")
	userID, _ := utils.GetUserID(c) // 这是一个难点，看下面说明

	// 1. 把用户的搜索词变成向量
	queryVec, err := h.ai.GetEmbedding(query)
	if err != nil {
		utils.Error(c, 500, "AI 服务繁忙")
		return
	}

	// 2. 去 Qdrant 搜出最相似的 Top 20 个 Note ID
	noteIDs, err := h.qdrant.Search(c, queryVec, 20, userID)
	if err != nil {
		utils.Error(c, 500, "搜索服务繁忙")
		return
	}

	if len(noteIDs) == 0 {
		utils.Success(c, []models.Note{})
		return
	}

	// 3. [关键] 回到 MySQL 查详情，并且加上 UserID 过滤！
	// 即使 Qdrant 搜到了别人的笔记，MySQL 这一步也会把它挡住，因为 Where 限制了 user_id
	// 这样 MySQL 就会正确返回这两种类型的笔记，同时依然拦截掉“别人的私有笔记”（如果有漏网之鱼的话）
	var notes []models.Note
	err = h.db.Where("id IN ?", noteIDs).
		Where(
			h.db.Where("user_id = ?", userID).
				Or("is_private = ?", false),
		).
		Find(&notes).Error

	// 4. (可选) 重新按照 Qdrant 返回的 ID 顺序排序 notes
	// 因为 MySQL 的 IN 查询返回顺序是不定的

	utils.Success(c, notes)
}
