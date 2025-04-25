package token

import (
    "errors"
    "time"
)

var (
    ErrInvalidToken = errors.New("token is invalid")
    ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
    UserID    uint      `json:"user_id"`
    IssuedAt  time.Time `json:"issued_at"`
    ExpiredAt time.Time `json:"expired_at"`
}

type Maker interface {
    CreateToken(userID uint, duration time.Duration) (string, error)
    
    VerifyToken(token string) (*Payload, error)
}