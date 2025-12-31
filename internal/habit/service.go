package habit

import "context"

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
