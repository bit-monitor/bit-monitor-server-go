package response

import "time"

type GetSubscriberNotifyRecord struct {
	Id         uint64    `json:"id"`
	State      int8      `json:"state"`
	Content    string    `json:"content"`
	Category   int8      `json:"category"`
	CreateTime time.Time `json:"createTime"`
}

type GetWithRelatedInfoSubscriberNotifyRecord struct {
	Id         uint64    `json:"id"`
	State      int8      `json:"state"`
	Content    string    `json:"content"`
	Category   int8      `json:"category"`
	AlarmName  string    `json:"alarmName"`
	CreateTime time.Time `json:"createTime"`
}
