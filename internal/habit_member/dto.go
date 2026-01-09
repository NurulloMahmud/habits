package habitmember

import "time"

type habitMemberCreateRequest struct {
	UserID  int64 `json:"user_id"`
	HabitID int64 `json:"habit_id"`
}

type habitOwner struct {
	ID        int64   `json:"user_id"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type userHabitsResponse struct {
	HabitID       int64      `json:"habit_id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	DailyCount    *int64     `json:"daily_count,omitempty"`
	DailyDuration *int64     `json:"daily_duration,omitempty"`
	PrivacyStatus string     `json:"privacy_status"`
	Identifier    *string    `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	Owner         habitOwner `json:"owner"`
}
