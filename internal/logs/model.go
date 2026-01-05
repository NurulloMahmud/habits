package logs

import "time"

type UserInfo struct {
	UserID int64  `bson:"user_id,omitempty"`
	IP     string `bson:"ip"`
}

type ActivityLog struct {
	User       UserInfo  `bson:"user"`
	Method     string    `bson:"method"`
	Endpoint   string    `bson:"endpoint"`
	Status     int       `bson:"status"`
	DurationMS int64     `bson:"duration_ms"`
	Error      *string   `bson:"error"`
	CreatedAt  time.Time `bson:"created_at"`
}
