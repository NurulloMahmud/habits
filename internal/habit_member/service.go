package habitmember

import (
	"context"
	"errors"

	cx "github.com/NurulloMahmud/habits/pkg/context"
)

var (
	errAlreadyMember = errors.New("You are already a member of this habit")
)

type Service struct {
	repo HabitMemberRepository
}

func NewService(hmRepo HabitMemberRepository) Service {
	return Service{repo: hmRepo}
}

func (s *Service) joinHabit(ctx context.Context, req habitMemberCreateRequest) (string, error) {
	privacyType, err := s.repo.getHabitPrivacyType(ctx, req.HabitID)
	if err != nil {
		return "", err
	}

	member, err := s.repo.isMember(ctx, req.HabitID, req.UserID)
	if err != nil {
		return "", err
	}

	if member {
		return "", errAlreadyMember
	}

	user, _ := ctx.Value("user").(*cx.User)

	if privacyType == "public" || user.UserRole == "admin" {
		err = s.repo.createHabitMember(ctx, req)
		if err != nil {
			return "", err
		}

		return "Member joined successfully", nil
	}

	err = s.repo.createjoinRequest(ctx, req)
	if err != nil {
		return "", err
	}

	return "Join request has been sent to habit owner", err
}
