package logic

import (
	"context"
	"errors"

	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"
)

var ErrForbidden = errors.New("forbidden")

type QueryTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTaskLogic {
	return &QueryTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTaskLogic) QueryTask(taskId string, userId int64, isAdmin bool) (resp *types.TaskResp, err error) {
	task, err := l.svcCtx.TaskModel.FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	if !isAdmin && task.UserId != userId {
		return nil, ErrForbidden
	}

	var audioUrl string
	if task.AudioUrl.Valid {
		audioUrl = task.AudioUrl.String
	}

	var errMsg string
	if task.ErrorMsg.Valid {
		errMsg = task.ErrorMsg.String
	}

	return &types.TaskResp{
		TaskId:   task.TaskId,
		Status:   task.Status,
		Progress: task.Progress,
		AudioUrl: audioUrl,
		Error:    errMsg,
	}, nil
}
