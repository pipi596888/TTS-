package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func ParseUserIDFromRequest(r *http.Request, secret string) (int64, bool, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, false, errors.New("missing auth header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return 0, false, errors.New("invalid auth header")
	}

	tokenStr := strings.TrimSpace(parts[1])
	if tokenStr == "" {
		return 0, false, errors.New("empty token")
	}

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return 0, false, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, errors.New("invalid claims")
	}

	raw := claims["userId"]
	switch v := raw.(type) {
	case float64:
		return int64(v), int64(v) == 1, nil
	case int64:
		return v, v == 1, nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, false, err
		}
		return id, id == 1, nil
	default:
		return 0, false, errors.New("missing userId")
	}
}
