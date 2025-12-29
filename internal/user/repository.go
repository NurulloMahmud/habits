package user

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	Create(ctx context.Context, u User) (*User, error)
	Get(ctx context.Context, id int64, email string) (*User, error)
	List()
	Update()
	Delete()
}

type postgresRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, u User) (*User, error) {
	query := `
	INSERT INTO users (email, password_hash, first_name, last_name, user_role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, user_role`

	err := r.db.QueryRowContext(ctx, query, u.Email, u.PasswordHash.hash, u.FirstName, u.LastName, u.UserRole).Scan(&u.ID, &u.UserRole)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *postgresRepo) Get(ctx context.Context, id int64, email string) (*User, error) {
	query := `
	SELECT 
		id, 
		email, 
		password_hash, 
		first_name, 
		last_name, 
		user_role, 
		is_active, 
		is_locked, 
		failed_attempts, 
		last_failed_login, 
		created_at
	FROM users
	WHERE id = $1 OR email = $2`

	var user User
	err := r.db.QueryRowContext(ctx, query, id, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash.hash,
		&user.FirstName,
		&user.LastName,
		&user.UserRole,
		&user.IsActive,
		&user.IsLocked,
		&user.FailedAttempts,
		&user.LastFailedLogin,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *postgresRepo) List() {
	return
}

func (r *postgresRepo) Update() {
	return
}

func (r *postgresRepo) Delete() {
	return
}
