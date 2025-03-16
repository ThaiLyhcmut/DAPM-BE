package helper

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID                   string
	jwt.RegisteredClaims // Thêm các trường chuẩn như exp, iat
}

type ContextKey string

const Auth ContextKey = "xxxyyyzzzkkk"

func CreateJWT(id string) string {
	claims := Claims{
		ID: id, // Sử dụng email ở đây
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Hết hạn sau 1 ngày
		}, // Thời gian hết hạn của token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(jwtSecret) // Tạo token
	if err != nil {
		fmt.Println("Error creating token:", err)
		return ""
	}
	return tokenString
}

func ParseJWT(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET") // Lấy secret từ .env

	// Parse token với claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra phương thức ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("phương thức ký không hợp lệ")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}
	// Kiểm tra nếu token hợp lệ và ép kiểu claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token không hợp lệ")
	}
}
