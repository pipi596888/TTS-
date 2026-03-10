package main

import (
	"database/sql"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type SystemStatsResp struct {
	Users              int64 `json:"users"`
	Voices             int64 `json:"voices"`
	Tasks              int64 `json:"tasks"`
	TasksPending       int64 `json:"tasksPending"`
	TasksProcessing    int64 `json:"tasksProcessing"`
	TasksSuccess       int64 `json:"tasksSuccess"`
	TasksFailed        int64 `json:"tasksFailed"`
	FeedbackOpen       int64 `json:"feedbackOpen"`
	CustomVoicePending int64 `json:"customVoicePending"`
}

func systemStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp SystemStatsResp

		_ = db.QueryRow("SELECT COUNT(1) FROM user").Scan(&resp.Users)
		_ = db.QueryRow("SELECT COUNT(1) FROM voice").Scan(&resp.Voices)
		_ = db.QueryRow("SELECT COUNT(1) FROM tts_task").Scan(&resp.Tasks)
		_ = db.QueryRow("SELECT COUNT(1) FROM tts_task WHERE status = 'pending'").Scan(&resp.TasksPending)
		_ = db.QueryRow("SELECT COUNT(1) FROM tts_task WHERE status = 'processing'").Scan(&resp.TasksProcessing)
		_ = db.QueryRow("SELECT COUNT(1) FROM tts_task WHERE status = 'success'").Scan(&resp.TasksSuccess)
		_ = db.QueryRow("SELECT COUNT(1) FROM tts_task WHERE status = 'failed'").Scan(&resp.TasksFailed)
		_ = db.QueryRow("SELECT COUNT(1) FROM feedback WHERE status = 'open'").Scan(&resp.FeedbackOpen)
		_ = db.QueryRow("SELECT COUNT(1) FROM custom_voice_request WHERE status = 'pending'").Scan(&resp.CustomVoicePending)

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
