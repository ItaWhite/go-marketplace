package security

import (
	"crypto/rsa"
	"fmt"
	"go-marketplace/internal/core/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator struct {
	publicKey *rsa.PublicKey
	issuer    string
}

func NewJWTValidator(issuer string, publicKey *rsa.PublicKey) *JWTValidator {
	return &JWTValidator{
		issuer:    issuer,
		publicKey: publicKey,
	}
}

type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func (v *JWTValidator) Validate(tokenString string) (AccessClaims, error) {
	var claims AccessClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}
			return v.publicKey, nil
		},
		jwt.WithIssuer(v.issuer),
	)
	if err != nil || !token.Valid {
		return AccessClaims{}, core_errors.ErrInvalidToken
	}

	return claims, nil
}
