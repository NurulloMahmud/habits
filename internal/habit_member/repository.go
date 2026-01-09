package habitmember

import (
	"context"
	"database/sql"
)

type HabitMemberRepository interface {
	createHabitMember(ctx context.Context, req habitMemberCreateRequest) error
	createjoinRequest(ctx context.Context, req habitMemberCreateRequest) error
	isMember(ctx context.Context, habitID, userID int64) (bool, error)
	getUserHabits(ctx context.Context, userID int64) ([]*userHabitsResponse, error)
	getHabitPrivacyType(ctx context.Context, habitID int64) (string, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) HabitMemberRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) createHabitMember(ctx context.Context, req habitMemberCreateRequest) error {
	query := `INSERT INTO habit_members (user_id, habit_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, req.UserID, req.HabitID)
	return err
}

func (r *postgresRepository) createjoinRequest(ctx context.Context, req habitMemberCreateRequest) error {
	query := `INSERT INTO habit_follow_requests (user_id, habit_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, req.UserID, req.HabitID)
	return err
}

func (r *postgresRepository) isMember(ctx context.Context, habitID, userID int64) (bool, error) {
	var result bool
	query := `SELECT EXISTS(SELECT 1 FROM habit_members WHERE habit_id = $1 AND user_id = $2)`
	err := r.db.QueryRowContext(ctx, query, habitID, userID).Scan(&result)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *postgresRepository) getUserHabits(ctx context.Context, userID int64) ([]*userHabitsResponse, error) {
	var result []*userHabitsResponse

	query := `
	SELECT 
		h.created_by,
		(SELECT email FROM users WHERE id = h.created_by) email,
		(SELECT first_name FROM users WHERE id = h.created_by) first_name,
		(SELECT last_name FROM users WHERE id = h.created_by) last_name,
		h.id,
		h.name,
		h.description,
		h.start_date,
		h.end_date,
		h.daily_count,
		h.daily_duration,
		h.privacy_status,
		h.identifier,
		h.created_at
	FROM 
		habit_members hms
		JOIN habits h ON h.id = hms.habit_id
	WHERE 
		hms.user_id = $1
	ORDER BY hms.id DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var owner habitOwner
		var userHabit userHabitsResponse

		err = rows.Scan(
			&owner.ID,
			&owner.Email,
			&owner.FirstName,
			&owner.LastName,
			&userHabit.HabitID,
			&userHabit.Name,
			&userHabit.Description,
			&userHabit.StartDate,
			&userHabit.EndDate,
			&userHabit.DailyCount,
			&userHabit.DailyDuration,
			&userHabit.PrivacyStatus,
			&userHabit.Identifier,
			&userHabit.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		userHabit.Owner = owner
		result = append(result, &userHabit)
	}

	return result, nil
}

func (r *postgresRepository) getHabitPrivacyType(ctx context.Context, habitID int64) (string, error) {
	var result string
	query := `SELECT privacy_status FROM habits WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, habitID).Scan(&result)
	return result, err
}
