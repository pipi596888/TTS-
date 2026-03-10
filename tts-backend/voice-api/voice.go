package main

import (
	"flag"
	"fmt"
	"net/http"

	"tts-backend/voice-api/internal/config"
	"tts-backend/voice-api/internal/handler"
	"tts-backend/voice-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/voice-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	svcCtx := svc.NewServiceContext(&c)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/voice/list",
				Handler: handler.GetVoiceListHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/voice/create",
				Handler: handler.CreateVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/voice/:id",
				Handler: handler.DeleteVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/voice/default/:id",
				Handler: handler.SetDefaultVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/voice/custom/request",
				Handler: handler.CreateCustomVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/voice/custom/list",
				Handler: handler.ListMyCustomVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/voice/custom/:id",
				Handler: handler.DeleteCustomVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/admin/voice/custom/list",
				Handler: handler.AdminListCustomVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/admin/voice/custom/:id/approve",
				Handler: handler.AdminApproveCustomVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/admin/voice/custom/:id/reject",
				Handler: handler.AdminRejectCustomVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/admin/voice/custom/:id",
				Handler: handler.AdminDeleteCustomVoiceHandler(svcCtx),
			},
		},
	)

	fmt.Printf("Starting voice-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
