package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const jwtIssuer = "weather-api"

type JWTManager struct {
	secret    []byte
	ttl       time.Duration
	blacklist TokenBlacklist
}

type TokenBlacklist interface {
	Revoke(ctx context.Context, tokenHash string, expiresAt time.Time) error
	Contains(ctx context.Context, tokenHash string) (bool, error)
}

func NewJWTManager(secret []byte, ttl time.Duration, blacklist TokenBlacklist) *JWTManager {
	return &JWTManager{secret: secret, ttl: ttl, blacklist: blacklist}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Generate(userID uuid.UUID, email string, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    jwtIssuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func (m *JWTManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func (m *JWTManager) Revoke(ctx context.Context, tokenString string) error {
	claims, err := m.Validate(tokenString)
	if err != nil {
		return err
	}

	if m.blacklist == nil {
		return nil
	}

	return m.blacklist.Revoke(ctx, hashToken(tokenString), claims.ExpiresAt.Time)
}

func (m *JWTManager) IsRevoked(ctx context.Context, tokenString string) (bool, error) {
	if m.blacklist == nil {
		return false, nil
	}

	return m.blacklist.Contains(ctx, hashToken(tokenString))
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
