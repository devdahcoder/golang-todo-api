package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Find(ctx context.Context, id uint) (*user.User, error) {
    query := `SELECT id, username, email, password, created_at, updated_at 
              FROM users WHERE id = $1`
              
    row := r.db.QueryRowContext(ctx, query, id)
    
    var u user.User
    err := row.Scan(
        &u.ID,
        &u.Username,
        &u.Email,
        &u.Password,
        &u.CreatedAt,
        &u.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // User not found, not an error
        }
        return nil, err
    }
    
    return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    return nil, nil
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
    query := `INSERT INTO users (username, email, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
              
    row := r.db.QueryRowContext(
        ctx,
        query,
        u.Username,
        u.Email,
        u.Password,
        u.CreatedAt,
        u.UpdatedAt,
    )
    
    return row.Scan(&u.ID)
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
    return nil, nil
}