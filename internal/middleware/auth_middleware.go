package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"note/config"
	"note/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func JWTAuthMiddleware(cfg *config.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		isBlacklisted, redisErr := utils.IsTokenBlacklisted(redisClient, tokenString)

		if redisErr != nil {
			// 💡 降级策略：Redis不可用时，跳过黑名单检查，只验证token签名
			slog.Warn("Redis unavailable, skipping blacklist check",
				"error", redisErr,
				"token", utils.GetTokenHash(tokenString))

		} else if isBlacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			return
		}

		token, err := utils.ValidateToken(cfg, tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, err := utils.ExtractClaims(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "extract claims failed"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Next()
	}
}
