package main

import (
	"database/sql"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type AdminRole struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type AdminRolesResp struct {
	List []AdminRole `json:"list"`
}

func listAdminRolesHandler(_ *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles := []AdminRole{
			{
				Key:         "admin",
				Name:        "管理员",
				Description: "拥有后台全部权限",
				Permissions: []string{"users:read", "users:write", "roles:read", "logs:read", "logs:write", "system:read"},
			},
			{
				Key:         "engineer",
				Name:        "工程师",
				Description: "可查看数据与日志（示例）",
				Permissions: []string{"logs:read", "system:read"},
			},
			{
				Key:         "user",
				Name:        "普通用户",
				Description: "仅能使用前台功能",
				Permissions: []string{},
			},
		}
		httpx.OkJsonCtx(r.Context(), w, AdminRolesResp{List: roles})
	}
}
