package note

import (
	"note/internal/ai"
	"note/internal/cache"
	"note/internal/mq"
	"note/internal/storage"
	"note/internal/vector"

	"gorm.io/gorm"
)

type NoteHandler struct {
	db             *gorm.DB
	cache          *cache.RedisCache
	rabbit         *mq.RabbitMQ
	ai             *ai.AIService
	qdrant         *vector.QdrantService
	storageService *storage.FileStorage
}

func NewNoteHandler(db *gorm.DB, cache *cache.RedisCache, rabbitMQ *mq.RabbitMQ, ai *ai.AIService, qdrant *vector.QdrantService, storageService *storage.FileStorage) *NoteHandler {
	return &NoteHandler{db: db, cache: cache, rabbit: rabbitMQ, ai: ai, qdrant: qdrant, storageService: storageService}
}
