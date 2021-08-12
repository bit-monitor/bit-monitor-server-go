package validation

type GetSubscriberNotifyRecord struct {
	PageInfo
	AlarmRecordId *uint64 `form:"alarmRecordId" json:"alarmRecordId" binding:"omitempty"`
	SubscriberId  *uint64 `form:"subscriberId" json:"subscriberId" binding:"omitempty"`
	State         *int8   `form:"state" json:"state" binding:"omitempty"`
	Content       string  `form:"content" json:"content" binding:"omitempty"`
	StartTime     string  `form:"startTime" json:"startTime" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime       string  `form:"endTime" json:"endTime" binding:"omitempty,datetime=2006-01-02 15:04:05"`
}
