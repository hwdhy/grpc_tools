package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

// CustomClaims 自定义结构
type CustomClaims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 根据用户ID生成token
func GenerateToken(userID uint64, role string, tokenKey string) string {
	claims := &CustomClaims{
		userID,
		role,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // expired 24 hour
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return ""
	}
	return tokenStr
}

// ParseToken 解析用户ID
func ParseToken(tokenStr string, tokenKey string) *jwt.Token {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		logrus.Printf("parse token error: %v", err)
		return nil
	}
	if token.Valid {
		return token
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			logrus.Println("not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			logrus.Println("token is expired or not valid!")
		} else {
			logrus.Println("couldn't handle this token:", err)
		}
	} else {
		logrus.Println("couldn't handle this token:", err)
	}
	return nil
}

// GetUserID 获取用户ID
func GetUserID(tokenStr string, tokenKey string) (uint64, string) {
	token := ParseToken(tokenStr, tokenKey)
	if token == nil {
		return 0, ""
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.UserID, claims.Role
	}
	return 0, ""
}
