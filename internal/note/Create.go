package note

import (
	"encoding/json"
	"fmt"
	"net/http"
	"note/internal/models"
	"note/internal/utils"
	"note/internal/validators"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	needSummary := c.DefaultQuery("gen_summary", "false") == "true"

	var req validators.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid note")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "生成中..."
	}

	var tags []models.Tag
	if len(req.TagIDs) > 0 {
		h.db.Where("id IN ?", req.TagIDs).Find(&tags)
	}

	note := models.Note{
		UserID:    userID,
		Title:     title,
		Content:   req.Content,
		Tags:      tags,
		IsPrivate: req.IsPrivate,
	}

	h.db.Create(&note)

	cacheKeyAllNotes := fmt.Sprintf("notes:user:%d", userID)
	h.cache.Del(c, cacheKeyAllNotes)

	go func() {
		// 场景 A: 如果标题为空（之前被处理成占位符了），发送生成标题任务
		if strings.TrimSpace(req.Title) == "" {
			h.sendAITask(note.ID, "generate_title")
		}

		// 场景 B: 如果前端要求生成摘要，发送生成摘要任务
		if needSummary {
			h.sendAITask(note.ID, "generate_summary")
		}
	}()

	if !note.IsPrivate {
		go func() {
			msg := models.FeedMsg{
				AuthorID: note.UserID,
				NoteID:   note.ID,
				PostTime: note.CreatedAt.Unix(),
			}
			body, _ := json.Marshal(msg)
			if h.rabbit != nil {
				// 只需要发这一条消息，剩下的交给消费者去扩散
				h.rabbit.Publish("feed_queue", body)
			}
		}()
	}
	utils.Success(c, note)
}

func (h *NoteHandler) sendAITask(noteID uint, taskType string) {
	if h.rabbit == nil {
		return
	}
	msg := models.AITaskMsg{
		NoteID: noteID,
		Task:   taskType,
	}
	body, _ := json.Marshal(msg)
	h.rabbit.Publish("ai_queue", body)
}
