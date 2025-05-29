package user

import (
	"context"
	"errors"
	"time"

	"github.com/devdahcoder/golang-todo-api/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)


type Service interface {
    GetUser(ctx context.Context, id uint) (*User, error)
    CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
    UpdateUser(ctx context.Context, id uint, input UpdateUserInput) (*User, error)
    DeleteUser(ctx context.Context, id uint) error
    Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
    ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}

type service struct {
	repo Repository
    tokenMaker token.Maker
}

func NewService(repo Repository, tokeMaker token.Maker) Service {
	return &service{
		repo: repo,
        tokenMaker: tokeMaker,
	}
}

func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
    user, err := s.repo.Find(ctx, id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Check if email already exists
    existingUser, err := s.repo.FindByEmail(ctx, input.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, ErrEmailAlreadyExists
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    now := time.Now()
    user := &User{
        Username:  input.Username,
        Email:     input.Email,
        Password:  string(hashedPassword),
        CreatedAt: now,
        UpdatedAt: now,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *service) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
    user, err := s.repo.FindByEmail(ctx, input.Email)
    if err != nil {
        return nil, err
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    token, err := s.tokenMaker.CreateToken(user.ID, 24*time.Hour)
    if err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        User:  *user,
        Token: token,
    }, nil
}

func (s *service) UpdateUser(ctx context.Context, id uint, input UpdateUserInput) (*User, error) {
    return nil, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
    return s.repo.List(ctx, limit, offset)
}