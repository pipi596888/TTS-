package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"tts-backend/voice-api/internal/model"
	"tts-backend/voice-api/internal/svc"
	"tts-backend/voice-api/internal/types"
)

var ErrForbidden = errors.New("forbidden")

type CustomVoiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCustomVoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CustomVoiceLogic {
	return &CustomVoiceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CustomVoiceLogic) Create(userId int64, req *types.CreateCustomVoiceReq) (*types.CreateCustomVoiceResp, error) {
	sampleUrlsJSON, _ := json.Marshal(req.SampleUrls)

	id, err := l.svcCtx.CustomVoiceRequestModel.Insert(&model.CustomVoiceRequest{
		UserId:     userId,
		Name:       req.Name,
		Tone:       req.Tone,
		Gender:     req.Gender,
		SampleText: req.SampleText,
		SampleUrls: string(sampleUrlsJSON),
		Status:     "pending",
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateCustomVoiceResp{Id: id}, nil
}

func (l *CustomVoiceLogic) ListForUser(userId int64) (*types.ListCustomVoiceResp, error) {
	items, err := l.svcCtx.CustomVoiceRequestModel.FindByUserId(userId, 50)
	if err != nil {
		return nil, err
	}

	resp := make([]types.CustomVoiceRequest, 0, len(items))
	for _, it := range items {
		resp = append(resp, toCustomVoiceResp(it))
	}
	return &types.ListCustomVoiceResp{List: resp, Total: int64(len(resp))}, nil
}

func (l *CustomVoiceLogic) ListAll(limit int) (*types.ListCustomVoiceResp, error) {
	items, err := l.svcCtx.CustomVoiceRequestModel.FindAll(limit)
	if err != nil {
		return nil, err
	}

	resp := make([]types.CustomVoiceRequest, 0, len(items))
	for _, it := range items {
		resp = append(resp, toCustomVoiceResp(it))
	}
	return &types.ListCustomVoiceResp{List: resp, Total: int64(len(resp))}, nil
}

func (l *CustomVoiceLogic) Delete(userId int64, id int64) error {
	// Ensure ownership (unless admin handled at handler layer)
	item, err := l.svcCtx.CustomVoiceRequestModel.FindOne(id)
	if err != nil {
		return err
	}
	if item.UserId != userId {
		return ErrForbidden
	}
	return l.svcCtx.CustomVoiceRequestModel.DeleteByIdAndUserId(id, userId)
}

func (l *CustomVoiceLogic) DeleteAsAdmin(id int64) error {
	return l.svcCtx.CustomVoiceRequestModel.Delete(id)
}

func (l *CustomVoiceLogic) Approve(id int64) (*types.ApproveCustomVoiceResp, error) {
	item, err := l.svcCtx.CustomVoiceRequestModel.FindOne(id)
	if err != nil {
		return nil, err
	}

	// Placeholder preview audio to make the approved voice immediately playable in UI.
	// In real workflows this should be a generated preview sample.
	const defaultPreviewURL = "https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav"

	voiceId, err := l.svcCtx.VoiceModel.Insert(&model.Voice{
		Name:       item.Name,
		Tone:       item.Tone,
		Gender:     item.Gender,
		PreviewUrl: defaultPreviewURL,
		IsDefault:  false,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.CustomVoiceRequestModel.UpdateApproval(
		id,
		"success",
		sql.NullInt64{Int64: voiceId, Valid: true},
		sql.NullString{Valid: false},
	)
	if err != nil {
		return nil, err
	}

	v, err := l.svcCtx.VoiceModel.FindOne(voiceId)
	if err != nil {
		return nil, err
	}

	return &types.ApproveCustomVoiceResp{
		Voice: types.Voice{
			Id:         v.Id,
			Name:       v.Name,
			Tone:       v.Tone,
			Gender:     v.Gender,
			PreviewUrl: v.PreviewUrl,
			IsDefault:  v.IsDefault,
		},
	}, nil
}

func (l *CustomVoiceLogic) Reject(id int64, errorMsg string) error {
	return l.svcCtx.CustomVoiceRequestModel.UpdateApproval(
		id,
		"failed",
		sql.NullInt64{Valid: false},
		sql.NullString{String: errorMsg, Valid: errorMsg != ""},
	)
}

func toCustomVoiceResp(it *model.CustomVoiceRequest) types.CustomVoiceRequest {
	var urls []string
	_ = json.Unmarshal([]byte(it.SampleUrls), &urls)

	r := types.CustomVoiceRequest{
		Id:         it.Id,
		UserId:     it.UserId,
		Name:       it.Name,
		Tone:       it.Tone,
		Gender:     it.Gender,
		SampleText: it.SampleText,
		SampleUrls: urls,
		Status:     it.Status,
		CreatedAt:  it.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  it.UpdatedAt.Format(time.RFC3339),
	}
	if it.ResultVoiceId.Valid {
		r.ResultVoiceId = it.ResultVoiceId.Int64
	}
	if it.ErrorMsg.Valid {
		r.ErrorMsg = it.ErrorMsg.String
	}
	return r
}
