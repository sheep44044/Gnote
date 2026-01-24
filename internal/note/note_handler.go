package note

import (
	"note/internal/ai"
	"note/internal/cache"
	"note/internal/mq"
	"note/internal/vector"

	"gorm.io/gorm"
)

type NoteHandler struct {
	db     *gorm.DB
	cache  *cache.RedisCache
	rabbit *mq.RabbitMQ
	ai     *ai.AIService
	qdrant *vector.QdrantService
}

func NewNoteHandler(db *gorm.DB, cache *cache.RedisCache, rabbitMQ *mq.RabbitMQ, ai *ai.AIService, qdrant *vector.QdrantService) *NoteHandler {
	return &NoteHandler{db: db, cache: cache, rabbit: rabbitMQ, ai: ai, qdrant: qdrant}
}
