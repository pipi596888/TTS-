package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tts-backend/voice-api/internal/auth"
	"tts-backend/voice-api/internal/logic"
	"tts-backend/voice-api/internal/svc"
	"tts-backend/voice-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}

		var req types.CreateCustomVoiceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "name is required"})
			return
		}
		if len(req.SampleUrls) == 0 {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "sampleUrls is required"})
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		resp, err := l.Create(userId, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func ListMyCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		resp, err := l.ListForUser(userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func DeleteCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}

		id, ok := parseIDWithPrefix(r.URL.Path, "/api/voice/custom/")
		if !ok {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "id is required"})
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		err = l.Delete(userId, id)
		if err != nil {
			if errors.Is(err, logic.ErrForbidden) {
				httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{"code": 403, "message": "forbidden"})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.Ok(w)
	}
}

func AdminListCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}
		if !isAdmin {
			httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{"code": 403, "message": "forbidden"})
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		resp, err := l.ListAll(200)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func AdminApproveCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}
		if !isAdmin {
			httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{"code": 403, "message": "forbidden"})
			return
		}

		id, ok := parseAdminCustomVoiceID(r.URL.Path, "/approve")
		if !ok {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "id is required"})
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		resp, err := l.Approve(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func AdminRejectCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}
		if !isAdmin {
			httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{"code": 403, "message": "forbidden"})
			return
		}

		id, ok := parseAdminCustomVoiceID(r.URL.Path, "/reject")
		if !ok {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "id is required"})
			return
		}

		var req types.RejectCustomVoiceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		if err := l.Reject(id, strings.TrimSpace(req.ErrorMsg)); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.Ok(w)
	}
}

func AdminDeleteCustomVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}
		if !isAdmin {
			httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{"code": 403, "message": "forbidden"})
			return
		}

		id, ok := parseAdminCustomVoiceID(r.URL.Path, "")
		if !ok {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "id is required"})
			return
		}

		l := logic.NewCustomVoiceLogic(r.Context(), svcCtx)
		if err := l.DeleteAsAdmin(id); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.Ok(w)
	}
}

func parseIDWithPrefix(path string, prefix string) (int64, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	s := strings.TrimPrefix(path, prefix)
	s = strings.Trim(s, "/")
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func parseAdminCustomVoiceID(path string, suffix string) (int64, bool) {
	// /api/admin/voice/custom/:id/approve
	const prefix = "/api/admin/voice/custom/"
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	s := strings.TrimPrefix(path, prefix)
	s = strings.Trim(s, "/")
	s = strings.TrimSuffix(s, suffix)
	s = strings.Trim(s, "/")
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
