package response

import "time"

type GetAlarm struct {
	Id                uint64 `json:"id"`
	Name              string `json:"name"`
	ProjectIdentifier string `json:"projectIdentifier"`
	Level             int8   `json:"level"`
	Category          int8   `json:"category"`
	Rule              string `json:"rule"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	SilentPeriod      int8   `json:"silentPeriod"`
	IsActive          int8   `json:"isActive"`
	CreateBy          uint64 `json:"createBy"`
	IsDeleted         int8   `json:"isDeleted"`
	SubscriberList    string `json:"subscriberList"`
}

type GetAlarmRecord struct {
	Id         uint64    `json:"id"`
	AlarmId    uint64    `json:"alarmId"`
	AlarmData  string    `json:"alarmData"`
	CreateTime time.Time `json:"createTime"`
	AlarmName  string    `json:"alarmName"`
}
