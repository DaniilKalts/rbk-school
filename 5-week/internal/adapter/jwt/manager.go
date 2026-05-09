package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	gojwt "github.com/golang-jwt/jwt/v5"
)

const issuer = "weather-api"

var ErrInvalidToken = errors.New("некорректный или просроченный токен")

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	gojwt.RegisteredClaims
}

type Blacklist interface {
	Revoke(ctx context.Context, token string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, token string) (bool, error)
}

type Manager struct {
	secret    []byte
	ttl       time.Duration
	blacklist Blacklist
}

func NewManager(secret []byte, ttl time.Duration, blacklist Blacklist) *Manager {
	return &Manager{secret: secret, ttl: ttl, blacklist: blacklist}
}

func (m *Manager) Generate(userID uuid.UUID, email string, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expiresAt),
			IssuedAt:  gojwt.NewNumericDate(now),
			Issuer:    issuer,
			Subject:   userID.String(),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func (m *Manager) Validate(tokenString string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenString, &Claims{}, func(token *gojwt.Token) (interface{}, error) {
		if token.Method != gojwt.SigningMethodHS256 {
			return nil, fmt.Errorf("неожиданный метод подписи: %s", token.Header["alg"])
		}

		return m.secret, nil
	}, gojwt.WithValidMethods([]string{gojwt.SigningMethodHS256.Alg()}), gojwt.WithExpirationRequired())

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *Manager) Revoke(ctx context.Context, tokenString string) error {
	claims, err := m.Validate(tokenString)
	if err != nil {
		return err
	}

	if m.blacklist == nil {
		return nil
	}

	return m.blacklist.Revoke(ctx, hashToken(tokenString), claims.ExpiresAt.Time)
}

func (m *Manager) IsRevoked(ctx context.Context, tokenString string) (bool, error) {
	if m.blacklist == nil {
		return false, nil
	}

	return m.blacklist.IsRevoked(ctx, hashToken(tokenString))
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
