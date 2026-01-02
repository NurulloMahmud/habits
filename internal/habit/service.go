package habit

import (
	"context"
	"errors"
	"strings"

	cx "github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/utils"
	"github.com/google/uuid"
)

var (
	errNoHabitFound = errors.New("No habit data found with given id/identifier")
	errNotOwner     = errors.New("You are not the owner of this habit")
	errTypeChange   = errors.New("Habit type cannot be changed. You can only update type's value")
	errHabitType    = errors.New("habit type can only be quantity or duration")
	errDateQuery    = errors.New("min date value must not be after max date")
	errInvalidSort  = errors.New("Invalid sort field value")
)

type Service struct {
	repo HabitRepository
}

func NewHabitService(repo HabitRepository) Service {
	return Service{repo: repo}
}

func (s *Service) create(ctx context.Context, data createHabitRequest) (*createHabitRequest, error) {
	habit, err := s.repo.create(ctx, data)
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *Service) update(ctx context.Context, userID int64, data updateHabitRequest) (*getHabitResponse, error) {
	habit, err := s.repo.get(ctx, int64(data.ID), "")
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errNoHabitFound
	}

	if userID != habit.Creator.ID {
		return nil, errNotOwner
	}

	// handle habit privacy type
	if data.PrivacyStatus != nil {
		if *data.PrivacyStatus != "public" && *data.PrivacyStatus != "private" {
			return nil, errInvalidStatus
		}
		if *data.PrivacyStatus == "private" && habit.PrivacyStatus == "public" {
			identifier := uuid.New().String()
			habit.Identifier = &identifier
			habit.PrivacyStatus = *data.PrivacyStatus
		}
		if *data.PrivacyStatus == "public" && habit.PrivacyStatus == "private" {
			habit.Identifier = nil
			habit.PrivacyStatus = *data.PrivacyStatus
		}
	}

	// do not allow habit type change (quantity based / duration based)
	if habit.DailyCount == nil && data.DailyCount != nil {
		return nil, errTypeChange
	}
	if habit.DailyDuration == nil && data.DailyDuration != nil {
		return nil, errTypeChange
	}

	if data.Name != nil {
		if strings.TrimSpace(*data.Name) == "" {
			return nil, errNameEmpty
		}
		habit.Name = *data.Name
	}
	if data.Description != nil {
		if strings.TrimSpace(*data.Description) == "" {
			return nil, errDescEmpty
		}
		habit.Description = *data.Description
	}

	if data.StartDate != nil || data.EndDate != nil {
		startDate := habit.StartDate
		endDate := habit.EndDate

		if data.StartDate != nil {
			startDate = *data.StartDate
		}
		if data.EndDate != nil {
			endDate = *data.EndDate
		}

		if !startDate.Before(endDate) {
			return nil, errInvalidDates
		}
	}

	if data.DailyCount != nil {
		habit.DailyCount = data.DailyCount
	}
	if data.DailyDuration != nil {
		habit.DailyDuration = data.DailyDuration
	}

	err = s.repo.update(ctx, *habit)
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *Service) delete(ctx context.Context, user cx.User, habitID int64) error {
	habit, err := s.repo.get(ctx, habitID, "")
	if err != nil {
		return err
	}
	if habit == nil {
		return errNoHabitFound
	}

	if user.ID != habit.Creator.ID && user.UserRole != "admin" {
		return errNotOwner
	}

	return s.repo.delete(ctx, habitID)
}

func (s *Service) list(ctx context.Context, query HabitListQuery) ([]*getHabitResponse, utils.Metadata, error) {
	var metaData utils.Metadata
	_, err := query.getHabitType()
	if err != nil {
		return nil, metaData, errHabitType
	}

	if query.startDate.minDate != nil && query.startDate.maxDate != nil {
		if query.startDate.maxDate.Before(*query.startDate.minDate) {
			return nil, metaData, errDateQuery
		}
	}
	if query.endDate.minDate != nil && query.endDate.maxDate != nil {
		if query.endDate.maxDate.Before(*query.endDate.minDate) {
			return nil, metaData, errDateQuery
		}
	}
	if query.createdAt.minDate != nil && query.createdAt.maxDate != nil {
		if query.createdAt.maxDate.Before(*query.createdAt.minDate) {
			return nil, metaData, errDateQuery
		}
	}

	data, metaData, err := s.repo.list(ctx, query)
	if err != nil {
		return nil, metaData, err
	}

	return data, metaData, nil
}
