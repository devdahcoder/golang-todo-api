package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
    secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
    if len(secretKey) < 32 {
        return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
    }
    return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userID uint, duration time.Duration) (string, error) {
    payload := &Payload{
        UserID:    userID,
        IssuedAt:  time.Now(),
        ExpiredAt: time.Now().Add(duration),
    }
    
    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":    payload.UserID,
        "issued_at":  payload.IssuedAt.Unix(),
        "expired_at": payload.ExpiredAt.Unix(),
    })
    
    return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
    keyFunc := func(token *jwt.Token) (interface{}, error) {
        _, ok := token.Method.(*jwt.SigningMethodHMAC)
        if !ok {
            return nil, ErrInvalidToken
        }
        return []byte(maker.secretKey), nil
    }
    
    jwtToken, err := jwt.Parse(token, keyFunc)
    if err != nil {
        return nil, err
    }
    
    claims, ok := jwtToken.Claims.(jwt.MapClaims)
    if !ok || !jwtToken.Valid {
        return nil, ErrInvalidToken
    }
    
    userID := uint(claims["user_id"].(float64))
    expiredAt := time.Unix(int64(claims["expired_at"].(float64)), 0)
    
    if time.Now().After(expiredAt) {
        return nil, ErrExpiredToken
    }
    
    payload := &Payload{
        UserID:    userID,
        IssuedAt:  time.Unix(int64(claims["issued_at"].(float64)), 0),
        ExpiredAt: expiredAt,
    }
    
    return payload, nil
}