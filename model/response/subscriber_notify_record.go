package response

import "time"

type GetSubscriberNotifyRecord struct {
	Id         uint64    `json:"id"`
	State      int8      `json:"state"`
	Content    string    `json:"content"`
	Category   int8      `json:"category"`
	CreateTime time.Time `json:"createTime"`
}
