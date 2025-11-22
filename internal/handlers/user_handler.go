package handlers

import (
	"net/http"
	"note/config"
	"note/internal/models"
	"note/internal/utils"
	"note/internal/validators"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{db: db, cfg: cfg}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req validators.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	var exists models.User
	if h.db.Where("username = ?", req.Username).First(&exists).RowsAffected > 0 {
		utils.Error(c, http.StatusConflict, "username already exists")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Username: req.Username,
		Password: string(hashed),
	}
	h.db.Create(&user)

	utils.Success(c, gin.H{"message": "user registered"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req validators.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	var user models.User
	if h.db.Where("username = ?", req.Username).First(&user).RowsAffected == 0 {
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	userIDStr := strconv.FormatUint(uint64(user.ID), 10)
	token, err := utils.GenerateToken(h.cfg, userIDStr, user.Username)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	utils.Success(c, gin.H{"token": token})
}
