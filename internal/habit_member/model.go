package habitmember

import (
	"time"

	"github.com/NurulloMahmud/habits/internal/habit"
	"github.com/NurulloMahmud/habits/internal/user"
)

type HabitMember struct {
	ID        int64       `json:"id"`
	Habit     habit.Habit `json:"habit"`
	Member    user.User   `json:"member"`
	CreatedAt time.Time   `json:"created_at"`
}

type HabitFollowRequest struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
