package habit

import (
	"context"
	"database/sql"
)

type HabitRepository interface {
	create(ctx context.Context, req createHabitRequest) (*createHabitRequest, error)
}

type postgresHabitRepository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) HabitRepository {
	return &postgresHabitRepository{db: db}
}

func (r *postgresHabitRepository) create(ctx context.Context, req createHabitRequest) (*createHabitRequest, error) {
	query := `
	INSERT INTO habits(name, description, start_date, end_date, daily_count, daily_duration, privacy_status, identifier, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6 * INTERVAL '1 minute', $7, $8, $9, $10)
	RETURNING id`

	err := r.db.QueryRowContext(
		ctx, query,
		req.Name,
		req.Description,
		req.StartDate,
		req.EndDate,
		req.DailyCount,
		req.DailyDuration,
		req.PrivacyStatus,
		req.Identifier,
		req.CreatedBy,
		req.CreatedAt,
	).Scan(&req.ID)

	if err != nil {
		return nil, err
	}

	return &req, nil
}
