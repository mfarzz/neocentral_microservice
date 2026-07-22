package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT signing configuration.
type JWTConfig struct {
	Secret             string
	AccessExpiry       time.Duration
	RefreshSecret      string
	RefreshExpiry      time.Duration
}

// Claims represents the JWT payload used by NeoCentral.
type Claims struct {
	jwt.RegisteredClaims
	Email   string   `json:"email,omitempty"`
	Roles   []string `json:"roles,omitempty"`
	Purpose string   `json:"purpose,omitempty"` // "access", "pwdreset", "verify"
}

// GenerateAccessToken signs an access token for the given user ID, email, and roles.
func GenerateAccessToken(userID, email string, roles []string, cfg JWTConfig) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessExpiry)),
		},
		Email: email,
		Roles: roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// GenerateRefreshToken signs a refresh token.
func GenerateRefreshToken(userID, email string, cfg JWTConfig) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.RefreshExpiry)),
		},
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.RefreshSecret))
}

// GeneratePurposeToken signs a special-purpose token (e.g. password reset, verification).
func GeneratePurposeToken(userID, purpose string, expiry time.Duration, cfg JWTConfig) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
		Purpose: purpose,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// VerifyAccessToken validates an access token and returns its claims.
func VerifyAccessToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// VerifyRefreshToken validates a refresh token.
func VerifyRefreshToken(tokenStr string, secret string) (*Claims, error) {
	return VerifyAccessToken(tokenStr, secret) // same logic, different secret
}
