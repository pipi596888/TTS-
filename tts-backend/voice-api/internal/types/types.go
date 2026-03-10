package types

type Voice struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Tone       string `json:"tone"`
	Gender     string `json:"gender"`
	PreviewUrl string `json:"previewUrl,omitempty"`
	IsDefault  bool   `json:"isDefault"`
}

type VoiceListReq struct{}

type VoiceListResp struct {
	List  []Voice `json:"list"`
	Total int64   `json:"total"`
}

type CreateVoiceReq struct {
	Name       string `json:"name"`
	Tone       string `json:"tone"`
	Gender     string `json:"gender"`
	PreviewUrl string `json:"previewUrl"`
}

type CreateVoiceResp struct {
	Voice Voice `json:"voice"`
}

type DeleteVoiceReq struct {
	Id int64 `json:"id"`
}

type SetDefaultReq struct {
	Id int64 `json:"id"`
}

type CustomVoiceRequest struct {
	Id            int64    `json:"id"`
	UserId        int64    `json:"userId"`
	Name          string   `json:"name"`
	Tone          string   `json:"tone"`
	Gender        string   `json:"gender"`
	SampleText    string   `json:"sampleText"`
	SampleUrls    []string `json:"sampleUrls"`
	Status        string   `json:"status"`
	ResultVoiceId int64    `json:"resultVoiceId,omitempty"`
	ErrorMsg      string   `json:"errorMsg,omitempty"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
}

type CreateCustomVoiceReq struct {
	Name       string   `json:"name"`
	Tone       string   `json:"tone"`
	Gender     string   `json:"gender"`
	SampleText string   `json:"sampleText"`
	SampleUrls []string `json:"sampleUrls"`
}

type CreateCustomVoiceResp struct {
	Id int64 `json:"id"`
}

type ListCustomVoiceReq struct{}

type ListCustomVoiceResp struct {
	List  []CustomVoiceRequest `json:"list"`
	Total int64                `json:"total"`
}

type DeleteCustomVoiceReq struct {
	Id int64 `json:"id"`
}

type ApproveCustomVoiceReq struct {
	Id int64 `json:"id"`
}

type ApproveCustomVoiceResp struct {
	Voice Voice `json:"voice"`
}

type RejectCustomVoiceReq struct {
	Id       int64  `json:"id"`
	ErrorMsg string `json:"errorMsg"`
}
