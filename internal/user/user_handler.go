package user

import (
	"note/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	rdb *redis.Client
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config, rdb *redis.Client) *UserHandler {
	return &UserHandler{db: db, cfg: cfg, rdb: rdb}
}
