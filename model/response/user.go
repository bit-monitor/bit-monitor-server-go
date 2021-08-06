package response

import "time"

type LoginUser struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	IsAdmin  int8   `json:"isAdmin"`
	Token    string `json:"token"`
}

type RelatedProject struct {
	Id                uint64    `json:"Id"`
	ProjectName       string    `json:"projectName"`
	ProjectIdentifier string    `json:"projectIdentifier"`
	Description       string    `json:"description"`
	AccessType        string    `json:"accessType"`
	ActiveFuncs       string    `json:"activeFuncs"`
	IsAutoUpload      uint8     `json:"isAutoUpload"`
	NotifyDtToken     string    `json:"notifyDtToken"`
	NotifyEmail       string    `json:"notifyEmail"`
	CreateTime        time.Time `json:"createTime"`
	UpdateTime        time.Time `json:"updateTime"`
}
