package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var JwtSecret string

func init() {
	JwtSecret = viper.GetString("JWT_SECRET")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(JwtSecret), nil
			})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		userID, _ := strconv.Atoi(claims.Subject)
		ctx := context.WithValue(r.Context(), "userID", int32(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWTToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwtSecret))
}
