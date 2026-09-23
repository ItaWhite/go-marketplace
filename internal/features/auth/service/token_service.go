package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	issuer     string
	privateKey *rsa.PrivateKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(issuer string, privateKey *rsa.PrivateKey, accessTTL, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		issuer:     issuer,
		privateKey: privateKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func (s *TokenService) GenerateAccessToken(userID int, role string) (string, error) {
	now := time.Now()

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(s.privateKey)
}

func (s *TokenService) GenerateRefreshToken() (string, time.Time, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)
	expiresAt := time.Now().Add(s.refreshTTL)

	return token, expiresAt, nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
