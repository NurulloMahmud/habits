package habitmember

import (
	"context"
	"database/sql"
)

type HabitMemberRepository interface {
	create(ctx context.Context, req habitMemberCreateRequest) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) HabitMemberRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) create(ctx context.Context, req habitMemberCreateRequest) error {
	query := `INSERT INTO habit_members (user_id, habit_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, req.UserID, req.HabitID)
	return err
}
