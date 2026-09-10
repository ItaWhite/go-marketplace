package security

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "user_role"
)

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

func UserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

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
		return AccessClaims{}, fmt.Errorf("invalid token")
	}

	return claims, nil
}
