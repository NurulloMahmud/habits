package habit

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NurulloMahmud/habits/pkg/utils"
)

type HabitRepository interface {
	create(ctx context.Context, req createHabitRequest) (*createHabitRequest, error)
	get(ctx context.Context, id int64, identifier string) (*getHabitResponse, error)
	update(ctx context.Context, data getHabitResponse) error
	delete(ctx context.Context, id int64) error
	list(ctx context.Context, q HabitListQuery) ([]*getHabitResponse, utils.Metadata, error)
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
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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

func (r *postgresHabitRepository) get(ctx context.Context, id int64, identifier string) (*getHabitResponse, error) {
	creator := habitCreator{}
	habit := getHabitResponse{}

	query := `
	SELECT 
		h.id id, 
		h.name name, 
		h.description deescription, 
		h.start_date start_date,  
		h.end_date end_date, 
		h.daily_count daily_count,  
		h.daily_duration daily_duration,
		h.privacy_status privacy_status,
		h.identifier identifier,
		h.created_at created_at,
		u.id creator_id,
		u.email creator_email,
		u.first_name creator_first_name,
		u.last_name creator_last_name
	FROM habits h
	JOIN users u ON u.id = h.created_by
	WHERE h.identifier = $1 OR h.id = $2`

	err := r.db.QueryRowContext(ctx, query, identifier, id).Scan(
		&habit.ID,
		&habit.Name,
		&habit.Description,
		&habit.StartDate,
		&habit.EndDate,
		&habit.DailyCount,
		&habit.DailyDuration,
		&habit.PrivacyStatus,
		&habit.Identifier,
		&habit.CreatedAt,
		&creator.ID,
		&creator.Email,
		&creator.FirstName,
		&creator.LastName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	habit.Creator = creator
	return &habit, nil
}

func (r *postgresHabitRepository) update(ctx context.Context, data getHabitResponse) error {
	query := `
	UPDATE habits
	SET name = $1, 
		description = $2, 
		start_date = $3, 
		end_date = $4, 
		daily_count = $5, 
		daily_duration = $6,
		privacy_status = $7,
		identifier = $8
	WHERE id = $9`

	_, err := r.db.ExecContext(
		ctx, query,
		data.Name,
		data.Description,
		data.StartDate,
		data.EndDate,
		data.DailyCount,
		data.DailyDuration,
		data.PrivacyStatus,
		data.Identifier,
		data.ID)

	return err
}

func (r *postgresHabitRepository) delete(ctx context.Context, id int64) error {
	query := `DELETE FROM habits WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresHabitRepository) list(ctx context.Context, q HabitListQuery) ([]*getHabitResponse, utils.Metadata, error) {
	var data []*getHabitResponse
	var metaData utils.Metadata
	var totalRecords int

	habitType, _ := q.getHabitType()
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) OVER(),
			h.id id, 
			h.name name, 
			h.description deescription, 
			h.start_date start_date,  
			h.end_date end_date, 
			h.daily_count daily_count,  
			h.daily_duration daily_duration,
			h.privacy_status privacy_status,
			h.identifier identifier,
			h.created_at created_at,
			u.id creator_id,
			u.email creator_email,
			u.first_name creator_first_name,
			u.last_name creator_last_name
		FROM habits h
		JOIN users u ON u.id = h.created_by
		WHERE 
			(h.name ILIKE $1 || '%%' OR $1 = '') AND
			(h.start_date >= $2 OR $2 IS NULL) AND
			(h.start_date <= $3 OR $3 IS NULL) AND
			(h.end_date <= $4 OR $4 IS NULL) AND
			(h.end_date <= $5 OR $5 IS NULL) AND
			(h.created_at <= $6 OR $6 IS NULL) AND
			(h.created_at <= $7 OR $7 IS NULL) AND
			(%s)
			ORDER BY %s, id
			LIMIT $8 OFFSET $9`, habitType, q.getSort())

	rows, err := r.db.QueryContext(
		ctx, query,
		q.search,
		q.startDate.minDate,
		q.startDate.maxDate,
		q.endDate.minDate,
		q.endDate.maxDate,
		q.createdAt.minDate,
		q.createdAt.maxDate,
		q.limit(),
		q.offset(),
	)

	if err != nil {
		return nil, metaData, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			creator habitCreator
			habit   getHabitResponse
		)
		err = rows.Scan(
			&totalRecords,
			&habit.ID,
			&habit.Name,
			&habit.Description,
			&habit.StartDate,
			&habit.EndDate,
			&habit.DailyCount,
			&habit.DailyDuration,
			&habit.PrivacyStatus,
			&habit.Identifier,
			&habit.CreatedAt,
			&creator.ID,
			&creator.Email,
			&creator.FirstName,
			&creator.LastName,
		)

		if err != nil {
			return nil, metaData, err
		}

		habit.Creator = creator
		data = append(data, &habit)
	}

	metaData = utils.CalculateMetadata(totalRecords, q.page, q.pageSize)
	return data, metaData, nil
}
