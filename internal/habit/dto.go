package habit

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	errNameEmpty            = errors.New("name field is required and cannot be empty")
	errDescEmpty            = errors.New("description field is required and cannot be empty")
	errStartDateEmpty       = errors.New("start_date field is required and cannot be empty")
	errEndDateEmpty         = errors.New("end_date field is required and cannot be empty")
	errTypeConflict         = errors.New("Habit must be either quantity based or duration. You cannot provide both daily_count and daily_duration fields")
	errTypeEmpty            = errors.New("Habit must be either quantity based or duration based. Please provide either daily_count or daily_duration")
	errInvalidStatus        = errors.New("privacy_status field must be either public or private and must not be empty")
	errInvalidDates         = errors.New("end_date cannot be before start_date")
	errInvalidDailyDuration = errors.New("daily duration minutes must be greater than or equal to 1")
)

type createHabitRequest struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	DailyCount    *int64     `json:"daily_count"`
	DailyDuration *int64     `json:"daily_duration"`
	PrivacyStatus string     `json:"privacy_status"`
	Identifier    *string    `json:"-"`
	CreatedBy     int64      `json:"-"`
	CreatedAt     time.Time  `json:"-"`
}

func (r *createHabitRequest) validateCreateRequest(createdBy int64) error {
	if strings.TrimSpace(r.Name) == "" {
		return errNameEmpty
	}
	if strings.TrimSpace(r.Description) == "" {
		return errDescEmpty
	}
	if r.StartDate == nil {
		return errStartDateEmpty
	}
	if r.EndDate == nil {
		return errEndDateEmpty
	}
	if r.DailyCount != nil && r.DailyDuration != nil {
		return errTypeConflict
	}
	if r.DailyCount == nil && r.DailyDuration == nil {
		return errTypeEmpty
	}
	if r.PrivacyStatus != "public" && r.PrivacyStatus != "private" {
		return errInvalidStatus
	}
	if r.DailyDuration != nil && *r.DailyDuration < 1 {
		return errInvalidDailyDuration
	}
	if r.EndDate.Before(*r.StartDate) {
		return errInvalidDates
	}

	if r.PrivacyStatus == "private" {
		identifier := uuid.New().String()
		r.Identifier = &identifier
	}

	r.CreatedBy = createdBy
	r.CreatedAt = time.Now().UTC()

	return nil
}

type updateHabitRequest struct {
	ID            int        `json:"-"`
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	DailyCount    *int64     `json:"daily_count"`
	DailyDuration *int64     `json:"daily_duration"`
	PrivacyStatus *string    `json:"privacy_status"`
	Identifier    *string    `json:"-"`
}

type habitCreator struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type getHabitResponse struct {
	ID            int64        `json:"id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	StartDate     time.Time    `json:"start_date"`
	EndDate       time.Time    `json:"end_date"`
	DailyCount    *int64       `json:"daily_count,omitempty"`
	DailyDuration *int64       `json:"daily_duration,omitempty"`
	PrivacyStatus string       `json:"privacy_status"`
	Identifier    *string      `json:"-"`
	CreatedAt     time.Time    `json:"created_at"`
	Creator       habitCreator `json:"creator"`
}

type dateFilter struct {
	minDate *time.Time
	maxDate *time.Time
}

type HabitListQuery struct {
	search      string
	habitType   string
	privacyType string
	startDate   dateFilter
	endDate     dateFilter
	createdAt   dateFilter
	pageSize    int
	page        int
	sort        string
	sortSafe    []string
	userRole    string
}

func (h *HabitListQuery) getHabitType() (string, error) {
	if h.habitType == "" {
		return " 1=1", nil
	}
	if h.habitType == "quantity" {
		return " h.daily_duration IS NULL ", nil
	}
	if h.habitType == "duration" {
		return " h.daily_count IS NULL ", nil
	}

	return "", errHabitType
}

func (h *HabitListQuery) validateSort() error {
	for _, val := range h.sortSafe {
		if (h.sort == val) || (string(h.sort[0]) == "-" && h.sort[1:] == val) {
			return nil
		}
	}
	return errInvalidSort
}

func (h *HabitListQuery) getSort() string {
	if string(h.sort[0]) == "-" {
		return h.sort[1:] + " DESC"
	}
	return h.sort
}

func (h *HabitListQuery) limit() int {
	return h.pageSize
}

func (h *HabitListQuery) offset() int {
	return h.page
}
