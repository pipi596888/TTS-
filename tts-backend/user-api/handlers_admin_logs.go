package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type AdminLogItem struct {
	Id            int64  `json:"id"`
	ActorUserId   int64  `json:"actorUserId"`
	ActorUsername string `json:"actorUsername"`
	Action        string `json:"action"`
	Ip            string `json:"ip"`
	CreatedAt     string `json:"createdAt"`
}

type AdminLogsResp struct {
	List  []AdminLogItem `json:"list"`
	Total int64          `json:"total"`
}

func listAdminLogsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize

		where := "1=1"
		args := make([]any, 0)
		if keyword != "" {
			where = "(l.action LIKE ? OR u.username LIKE ?)"
			kw := "%" + keyword + "%"
			args = append(args, kw, kw)
		}

		var total int64
		countSQL := fmt.Sprintf("SELECT COUNT(1) FROM admin_log l LEFT JOIN user u ON u.id = l.actor_user_id WHERE %s", where)
		if err := db.QueryRow(countSQL, args...).Scan(&total); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		querySQL := fmt.Sprintf(`
SELECT l.id, l.actor_user_id, COALESCE(u.username,''), l.action, COALESCE(l.ip,'-'), l.created_at
FROM admin_log l
LEFT JOIN user u ON u.id = l.actor_user_id
WHERE %s
ORDER BY l.id DESC
LIMIT ? OFFSET ?`, where)
		args = append(args, pageSize, offset)

		rows, err := db.Query(querySQL, args...)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer rows.Close()

		list := make([]AdminLogItem, 0)
		for rows.Next() {
			var item AdminLogItem
			var createdAt []byte
			if err := rows.Scan(&item.Id, &item.ActorUserId, &item.ActorUsername, &item.Action, &item.Ip, &createdAt); err != nil {
				continue
			}
			item.CreatedAt = string(createdAt)
			list = append(list, item)
		}

		httpx.OkJsonCtx(r.Context(), w, AdminLogsResp{
			List:  list,
			Total: total,
		})
	}
}

type AdminAppendLogReq struct {
	Action string `json:"action"`
}

func appendAdminLogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AdminAppendLogReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid body")
			return
		}
		action := strings.TrimSpace(req.Action)
		if action == "" {
			writeJSONError(w, http.StatusBadRequest, "action is required")
			return
		}

		actor := mustUserID(r.Context())
		writeAdminLog(db, actor, action, getClientIP(r))

		httpx.OkJsonCtx(r.Context(), w, map[string]any{"ok": true})
	}
}
