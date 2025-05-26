package login

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	"github.com/kp/pager/login/models"
)

var jwtKey = []byte("your-secret-key") // TODO: Move to config

// InitCache should be called during application startup
// with the Redis server address

type Claims struct {
	Username    string   `json:"username"`
	UserType    string   `json:"user_type"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

func GenerateToken(user *models.User, permissions []string) (string, time.Time, error) {
	expirationTime := time.Now().Add(72 * time.Hour)

	claims := &Claims{
		Username:    user.Username,
		UserType:    user.UserType,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, expirationTime, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func AuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("X-Auth-Token")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Authorization header required"}`))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ValidateToken(tokenString)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Invalid token"}`))
				return
			}

			// Add claims to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "user_type", claims.UserType)
			ctx = context.WithValue(ctx, "permissions", claims.Permissions)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HasPermission(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get permissions from token
			permissions, ok := r.Context().Value("permissions").([]string)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "No permissions found"}`))
				return
			}

			// Check if token has required permission
			for _, perm := range permissions {
				if perm == requiredPermission {
					// Update cache
					CacheUserPermission(r.Context(), requiredPermission)
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Insufficient permissions"}`))
		})
	}
}
