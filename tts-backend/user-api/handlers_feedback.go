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

type Feedback struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"userId"`
	Username  string `json:"username,omitempty"`
	Category  string `json:"category"`
	Content   string `json:"content"`
	Contact   string `json:"contact"`
	Status    string `json:"status"`
	Reply     string `json:"reply,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateFeedbackReq struct {
	Category string `json:"category"`
	Content  string `json:"content"`
	Contact  string `json:"contact"`
}

type CreateFeedbackResp struct {
	Id int64 `json:"id"`
}

type FeedbackListResp struct {
	List  []Feedback `json:"list"`
	Total int64      `json:"total"`
}

type ReplyFeedbackReq struct {
	Reply  string `json:"reply"`
	Status string `json:"status"`
}

func createFeedbackHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())

		var req CreateFeedbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "参数错误")
			return
		}
		req.Content = strings.TrimSpace(req.Content)
		if req.Content == "" {
			writeJSONError(w, http.StatusBadRequest, "反馈内容不能为空")
			return
		}

		res, err := db.Exec(
			"INSERT INTO feedback (user_id, category, content, contact, status) VALUES (?, ?, ?, ?, 'open')",
			userId,
			strings.TrimSpace(req.Category),
			req.Content,
			strings.TrimSpace(req.Contact),
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id, _ := res.LastInsertId()

		httpx.OkJsonCtx(r.Context(), w, CreateFeedbackResp{Id: id})
	}
}

func listMyFeedbackHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())

		rows, err := db.Query(
			"SELECT id, user_id, COALESCE(category,''), content, COALESCE(contact,''), status, COALESCE(reply,''), created_at, updated_at FROM feedback WHERE user_id = ? ORDER BY created_at DESC LIMIT 50",
			userId,
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer rows.Close()

		list := make([]Feedback, 0)
		for rows.Next() {
			var f Feedback
			var createdAt, updatedAt []byte
			if err := rows.Scan(
				&f.Id, &f.UserId, &f.Category, &f.Content, &f.Contact, &f.Status, &f.Reply, &createdAt, &updatedAt,
			); err != nil {
				continue
			}
			f.CreatedAt = string(createdAt)
			f.UpdatedAt = string(updatedAt)
			list = append(list, f)
		}

		httpx.OkJsonCtx(r.Context(), w, FeedbackListResp{List: list, Total: int64(len(list))})
	}
}

func listAllFeedbackHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(
			`SELECT f.id, f.user_id, u.username, COALESCE(f.category,''), f.content, COALESCE(f.contact,''), f.status, COALESCE(f.reply,''), f.created_at, f.updated_at
			 FROM feedback f
			 LEFT JOIN user u ON u.id = f.user_id
			 ORDER BY f.created_at DESC LIMIT 200`,
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer rows.Close()

		list := make([]Feedback, 0)
		for rows.Next() {
			var f Feedback
			var createdAt, updatedAt []byte
			if err := rows.Scan(
				&f.Id, &f.UserId, &f.Username, &f.Category, &f.Content, &f.Contact, &f.Status, &f.Reply, &createdAt, &updatedAt,
			); err != nil {
				continue
			}
			f.CreatedAt = string(createdAt)
			f.UpdatedAt = string(updatedAt)
			list = append(list, f)
		}
		httpx.OkJsonCtx(r.Context(), w, FeedbackListResp{List: list, Total: int64(len(list))})
	}
}

func replyFeedbackHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /api/admin/feedback/:id/reply
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/feedback/")
		path = strings.Trim(path, "/")
		path = strings.TrimSuffix(path, "/reply")
		path = strings.Trim(path, "/")
		id, err := strconv.ParseInt(path, 10, 64)
		if err != nil || id <= 0 {
			writeJSONError(w, http.StatusBadRequest, "id is required")
			return
		}

		var req ReplyFeedbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "参数错误")
			return
		}
		req.Reply = strings.TrimSpace(req.Reply)
		if req.Reply == "" {
			writeJSONError(w, http.StatusBadRequest, "reply 不能为空")
			return
		}
		status := strings.TrimSpace(req.Status)
		if status == "" {
			status = "closed"
		}

		_, err = db.Exec("UPDATE feedback SET reply = ?, status = ? WHERE id = ?", req.Reply, status, id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		writeAdminLog(db, mustUserID(r.Context()), fmt.Sprintf("回复了反馈 #%d", id), getClientIP(r))

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"code":    0,
			"message": "ok",
		})
	}
}
