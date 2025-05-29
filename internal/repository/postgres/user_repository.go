package postgres

import (
	"context"
	// "database/sql"
	// "errors"

	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Find(ctx context.Context, id uint) (*user.User, error) {
    // query := `SELECT id, username, email, password, created_at, updated_at 
    //           FROM users WHERE id = $1`
              
    // row := r.db.QueryRowContext(ctx, query, id)
    
    // var u user.User
    // err := row.Scan(
    //     &u.ID,
    //     &u.Username,
    //     &u.Email,
    //     &u.Password,
    //     &u.CreatedAt,
    //     &u.UpdatedAt,
    // )
    
    // if err != nil {
    //     if errors.Is(err, sql.ErrNoRows) {
    //         return nil, nil // User not found, not an error
    //     }
    //     return nil, err
    // }
    
    // return &u, nil
    return nil, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    return nil, nil
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
    // query := `INSERT INTO users (username, email, password, created_at, updated_at) 
    //           VALUES ($1, $2, $3, $4, $5) RETURNING id`
              
    // row := r.db.QueryRowContext(
    //     ctx,
    //     query,
    //     u.Username,
    //     u.Email,
    //     u.Password,
    //     u.CreatedAt,
    //     u.UpdatedAt,
    // )
    
    // return row.Scan(&u.ID)
    return nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
    rows, err := r.db.Query(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users LIMIT $1 OFFSET $2", limit, offset)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var users []*user.User

    for rows.Next() {
        var u user.User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
            return nil, err
        }
        users = append(users, &u)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil

}