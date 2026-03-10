package types

type Segment struct {
	VoiceId int64  `json:"voiceId"`
	Emotion string `json:"emotion"`
	Text    string `json:"text"`
}

type GenerateReq struct {
	Segments []Segment `json:"segments"`
	Format   string    `json:"format"`
	Channel  string    `json:"channel"`
}

type GenerateResp struct {
	TaskId string `json:"taskId"`
}

type TaskResp struct {
	TaskId   string `json:"taskId"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	AudioUrl string `json:"audioUrl,omitempty"`
	Error    string `json:"error,omitempty"`
}

type TaskDetailSegment struct {
	VoiceId int64  `json:"voiceId"`
	Emotion string `json:"emotion"`
	Text    string `json:"text"`
	Sort    int    `json:"sort"`
}

type TaskDetailResp struct {
	TaskId    string              `json:"taskId"`
	Title     string              `json:"title"`
	Status    string              `json:"status"`
	Progress  int                 `json:"progress"`
	AudioUrl  string              `json:"audioUrl,omitempty"`
	Error     string              `json:"error,omitempty"`
	Format    string              `json:"format"`
	Channel   string              `json:"channel"`
	Segments  []TaskDetailSegment `json:"segments"`
	CreatedAt string              `json:"createdAt,omitempty"`
	UpdatedAt string              `json:"updatedAt,omitempty"`
}
