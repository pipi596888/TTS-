package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"tts-backend/user-api/internal/config"
)

type ctxKey string

const userIDKey ctxKey = "userId"

func requireAuth(c *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := parseUserIDFromRequest(r, c.JwtSecret)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "未登录或登录已过期")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userId)
		next(w, r.WithContext(ctx))
	}
}

func requireAdmin(c *config.Config, db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return requireAuth(c, func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())
		if !isAdminUserID(userId) {
			writeJSONError(w, http.StatusForbidden, "无权限")
			return
		}
		_ = db
		next(w, r)
	})
}

func mustUserID(ctx context.Context) int64 {
	v := ctx.Value(userIDKey)
	if v == nil {
		return 0
	}
	id, _ := v.(int64)
	return id
}

func isAdminUserID(userId int64) bool {
	// Simplified: treat the very first user as admin.
	return userId == 1
}

func parseUserIDFromRequest(r *http.Request, secret string) (int64, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return 0, errors.New("missing auth header")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return 0, errors.New("invalid auth header")
	}
	tokenStr := strings.TrimSpace(parts[1])
	if tokenStr == "" {
		return 0, errors.New("empty token")
	}

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	raw := claims["userId"]
	switch v := raw.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return id, nil
	default:
		return 0, errors.New("missing userId")
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    status,
		"message": message,
	})
}

func parseIDFromPath(prefix string, path string) (int64, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	s := strings.TrimPrefix(path, prefix)
	s = strings.Trim(s, "/")
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
