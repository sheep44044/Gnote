package utils

import (
	"context"
	"math/rand"
	"note/config"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func GenerateToken(cfg *config.Config, userID string, username string) (string, error) {
	// 生成唯一ID用于黑名单
	jti := time.Now().UnixNano() + rand.Int63()

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"jti":      jti,
		"exp":      time.Now().Add(cfg.JWTExpirationTime).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      cfg.JWTIssuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecretKey))
}

// 检查token是否在黑名单中
func IsTokenBlacklisted(redisClient *redis.Client, tokenString string) bool {
	// 先简单解析token获取jti，不验证签名（因为要先检查黑名单）
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false
	}

	// 只解析claims部分
	claims := jwt.MapClaims{}
	_, _, _ = jwt.NewParser().ParseUnverified(tokenString, claims)

	if jti, ok := claims["jti"].(float64); ok {
		key := "blacklist:" + strconv.FormatInt(int64(jti), 10)
		_, err := redisClient.Get(context.Background(), key).Result()
		return err == nil // 存在即被加入黑名单
	}
	return false
}

// 将token加入黑名单
func AddTokenToBlacklist(redisClient *redis.Client, tokenString string, expiration time.Duration) error {
	claims := jwt.MapClaims{}
	_, _, _ = jwt.NewParser().ParseUnverified(tokenString, claims)

	if jti, ok := claims["jti"].(float64); ok {
		key := "blacklist:" + strconv.FormatInt(int64(jti), 10)
		return redisClient.Set(context.Background(), key, "1", expiration).Err()
	}
	return nil
}

func ValidateToken(cfg *config.Config, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWTSecretKey), nil
	})
}

func ExtractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
