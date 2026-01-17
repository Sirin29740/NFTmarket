package api

import (
	"NFTmarket/internal/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		//if authHeader == "" {
		//	c.JSON(401, gin.H{"error": "missing authorization header"})
		//	c.Abort()
		//	return
		//}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token 格式错误，应为 Bearer Token"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			// 🚨 关键：在这里打印详细错误，确定是 'signature is invalid' 还是 'token is expired'
			fmt.Println("CRITICAL ERROR:", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "token无效或过期", "detail": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("suername", claims.Username)

		c.Next()
	}
}
