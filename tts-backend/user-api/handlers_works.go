package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Work struct {
	TaskId    string `json:"taskId"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Format    string `json:"format"`
	AudioUrl  string `json:"audioUrl"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type WorksResp struct {
	List  []Work `json:"list"`
	Total int64  `json:"total"`
}

func getWorksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())

		rows, err := db.Query(
			"SELECT task_id, COALESCE(title,''), status, progress, format, COALESCE(audio_url,''), COALESCE(error_msg,''), created_at FROM tts_task WHERE user_id = ? ORDER BY created_at DESC LIMIT 50",
			userId,
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer rows.Close()

		works := make([]Work, 0)
		for rows.Next() {
			var work Work
			var createdAt []byte
			if err := rows.Scan(
				&work.TaskId,
				&work.Title,
				&work.Status,
				&work.Progress,
				&work.Format,
				&work.AudioUrl,
				&work.ErrorMsg,
				&createdAt,
			); err != nil {
				continue
			}
			work.CreatedAt = string(createdAt)
			works = append(works, work)
		}

		httpx.OkJsonCtx(r.Context(), w, WorksResp{
			List:  works,
			Total: int64(len(works)),
		})
	}
}

type UpdateWorkTitleReq struct {
	Title string `json:"title"`
}

func updateWorkTitleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())

		taskId := strings.TrimPrefix(r.URL.Path, "/api/works/")
		taskId = strings.Trim(taskId, "/")
		taskId = strings.TrimSuffix(taskId, "/title")
		taskId = strings.Trim(taskId, "/")
		if taskId == "" {
			writeJSONError(w, http.StatusBadRequest, "taskId is required")
			return
		}

		var req UpdateWorkTitleReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid body")
			return
		}
		title := strings.TrimSpace(req.Title)
		if title == "" {
			writeJSONError(w, http.StatusBadRequest, "title is required")
			return
		}
		if len(title) > 255 {
			title = title[:255]
		}

		res, err := db.Exec("UPDATE tts_task SET title = ? WHERE task_id = ? AND user_id = ?", title, taskId, userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			writeJSONError(w, http.StatusNotFound, "not found")
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"code":    0,
			"message": "ok",
		})
	}
}

func deleteWorkHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mustUserID(r.Context())

		taskId := strings.TrimPrefix(r.URL.Path, "/api/works/")
		taskId = strings.Trim(taskId, "/")
		if taskId == "" {
			writeJSONError(w, http.StatusBadRequest, "taskId is required")
			return
		}

		res, err := db.Exec("DELETE FROM tts_task WHERE task_id = ? AND user_id = ?", taskId, userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			writeJSONError(w, http.StatusNotFound, "not found")
			return
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"code":    0,
			"message": "ok",
		})
	}
}

func ensureWorksSchema(db *sql.DB) error {
	// Ensure `title` exists on `tts_task` to support rename without manual migrations.
	var count int
	err := db.QueryRow(
		"SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tts_task' AND COLUMN_NAME = 'title'",
	).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = db.Exec("ALTER TABLE tts_task ADD COLUMN title VARCHAR(255) DEFAULT '' AFTER task_id")
	return err
}
