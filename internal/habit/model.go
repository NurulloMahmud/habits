package habit

import (
	"time"
)

type Habit struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	DailyCount    *int64    `json:"daily_count"`
	DailyDuration *int64    `json:"daily_duration"`
	PrivacyStatus string    `json:"privacy_status"`
	Identifier    *string   `json:"-"`
	CreatedBy     int64     `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}
