package habitmember

type habitMemberCreateRequest struct {
	UserID  int64 `json:"user_id"`
	HabitID int64 `json:"habit_id"`
}
