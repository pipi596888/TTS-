package logic

import (
	"context"
	"time"

	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"
)

type QueryTaskDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTaskDetailLogic {
	return &QueryTaskDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTaskDetailLogic) QueryTaskDetail(taskId string, userId int64, isAdmin bool) (resp *types.TaskDetailResp, err error) {
	task, err := l.svcCtx.TaskModel.FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}
	if !isAdmin && task.UserId != userId {
		return nil, ErrForbidden
	}

	segments, err := l.svcCtx.SegmentModel.FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	outSegs := make([]types.TaskDetailSegment, 0, len(segments))
	for _, s := range segments {
		outSegs = append(outSegs, types.TaskDetailSegment{
			VoiceId: s.VoiceId,
			Emotion: s.Emotion,
			Text:    s.Text,
			Sort:    s.Sort,
		})
	}

	var audioUrl string
	if task.AudioUrl.Valid {
		audioUrl = task.AudioUrl.String
	}

	var errMsg string
	if task.ErrorMsg.Valid {
		errMsg = task.ErrorMsg.String
	}

	return &types.TaskDetailResp{
		TaskId:    task.TaskId,
		Title:     task.Title,
		Status:    task.Status,
		Progress:  task.Progress,
		AudioUrl:  audioUrl,
		Error:     errMsg,
		Format:    task.Format,
		Channel:   task.Channel,
		Segments:  outSegs,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}, nil
}
