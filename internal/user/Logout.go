package user

import (
	"net/http"
	"note/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Error(c, http.StatusBadRequest, "missing token")
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.Error(c, http.StatusBadRequest, "invalid token format")
		return
	}
	tokenString := parts[1]

	// 获取token剩余有效期 - 简单做法：使用配置的过期时间
	expiration := h.cfg.JWTExpirationTime

	// 加入黑名单
	err := utils.AddTokenToBlacklist(tokenString, expiration)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to logout")
		return
	}

	utils.Success(c, gin.H{"message": "logged out successfully"})
}
