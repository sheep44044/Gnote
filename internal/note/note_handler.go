package note

import (
	"note/internal/cache"
	"note/internal/mq"

	"gorm.io/gorm"
)

type NoteHandler struct {
	db     *gorm.DB
	cache  *cache.RedisCache
	rabbit *mq.RabbitMQ
}

func NewNoteHandler(db *gorm.DB, cache *cache.RedisCache, rabbitMQ *mq.RabbitMQ) *NoteHandler {
	return &NoteHandler{db: db, cache: cache, rabbit: rabbitMQ}
}
