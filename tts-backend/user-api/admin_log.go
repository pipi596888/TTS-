package main

import (
	"database/sql"
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}
	xri := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func writeAdminLog(db *sql.DB, actorUserId int64, action string, ip string) {
	action = strings.TrimSpace(action)
	if action == "" {
		return
	}
	if ip == "" {
		ip = "-"
	}
	_, _ = db.Exec("INSERT INTO admin_log (actor_user_id, action, ip) VALUES (?, ?, ?)", actorUserId, action, ip)
}
