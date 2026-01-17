package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("wcccPU5Fv752liiDHIGLDFesaP5V0EmhaIPfveQcrPE=")

type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	expiretime := time.Now().Add(24 * time.Hour)
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiretime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ValidateToken(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		// 🚨 关键：返回详细的错误信息
		// 打印到 Go 控制台：
		fmt.Println("JWT VALIDATION ERROR:", err.Error())
		// 返回给调用者，供前端查看
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}
	return claims, nil
}
