package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	core_errors "go-marketplace/internal/core/errors"
	"testing"
	"time"

	"go-marketplace/internal/core/domain"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type tokenRepositoryMock struct {
	createRefreshTokenFunc func(context.Context, domain.RefreshToken) error
	getRefreshTokenFunc    func(context.Context, string) (domain.RefreshToken, error)
	rotateRefreshTokenFunc func(context.Context, int, domain.RefreshToken) error
	revokeRefreshTokenFunc func(context.Context, int) error

	createdRefreshToken *domain.RefreshToken
	rotatedOldTokenID   int
	rotatedNewToken     *domain.RefreshToken
	revokedTokenID      int
}

func (m *tokenRepositoryMock) CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	if m.createdRefreshToken != nil {
		*m.createdRefreshToken = token
	}

	if m.createRefreshTokenFunc != nil {
		return m.createRefreshTokenFunc(ctx, token)
	}

	return nil
}

func (m *tokenRepositoryMock) GetRefreshToken(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	if m.getRefreshTokenFunc != nil {
		return m.getRefreshTokenFunc(ctx, tokenHash)
	}

	return domain.RefreshToken{}, nil
}

func (m *tokenRepositoryMock) RotateRefreshToken(ctx context.Context, oldTokenID int, newToken domain.RefreshToken) error {
	m.rotatedOldTokenID = oldTokenID

	if m.rotatedNewToken != nil {
		*m.rotatedNewToken = newToken
	}

	if m.rotateRefreshTokenFunc != nil {
		return m.rotateRefreshTokenFunc(ctx, oldTokenID, newToken)
	}

	return nil
}

func (m *tokenRepositoryMock) RevokeRefreshToken(ctx context.Context, id int) error {
	m.revokedTokenID = id

	if m.revokeRefreshTokenFunc != nil {
		return m.revokeRefreshTokenFunc(ctx, id)
	}

	return nil
}

type credentialsRepositoryMock struct {
	createCredentialsFunc func(context.Context, domain.Credentials) error
	getByLoginFunc        func(context.Context, string) (domain.Credentials, error)

	createdCredentials *domain.Credentials
}

func (m *credentialsRepositoryMock) CreateCredentials(ctx context.Context, credentials domain.Credentials) error {
	if m.createdCredentials != nil {
		*m.createdCredentials = credentials
	}

	if m.createCredentialsFunc != nil {
		return m.createCredentialsFunc(ctx, credentials)
	}

	return nil
}

func (m *credentialsRepositoryMock) GetByLogin(ctx context.Context, login string) (domain.Credentials, error) {
	if m.getByLoginFunc != nil {
		return m.getByLoginFunc(ctx, login)
	}

	return domain.Credentials{}, nil
}

type userServiceMock struct {
	createFunc  func(context.Context, CreateUserInput) (int, string, error)
	getRoleFunc func(context.Context, int) (string, error)

	createdInput *CreateUserInput
}

func (m *userServiceMock) Create(ctx context.Context, input CreateUserInput) (int, string, error) {
	if m.createdInput != nil {
		*m.createdInput = input
	}

	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}

	return 0, "", nil
}

func (m *userServiceMock) GetRole(ctx context.Context, id int) (string, error) {
	if m.getRoleFunc != nil {
		return m.getRoleFunc(ctx, id)
	}

	return "", nil
}

func newTestTokenService(t *testing.T) *TokenService {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	return NewTokenService(
		"test",
		key,
		15*time.Minute,
		7*24*time.Hour,
	)
}

func TestAuthService_Login(t *testing.T) {
	password := "password123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)

	tokenRepo := &tokenRepositoryMock{}
	credRepo := &credentialsRepositoryMock{
		getByLoginFunc: func(
			_ context.Context,
			login string,
		) (domain.Credentials, error) {
			require.Equal(t, "user", login)

			return domain.Credentials{
				UserID:       42,
				Login:        "user",
				PasswordHash: string(hash),
			}, nil
		},
	}

	userRepo := &userServiceMock{
		getRoleFunc: func(
			_ context.Context,
			id int,
		) (string, error) {
			require.Equal(t, 42, id)
			return "buyer", nil
		},
	}

	var createdToken domain.RefreshToken
	tokenRepo.createdRefreshToken = &createdToken

	service := NewAuthService(
		tokenRepo,
		credRepo,
		userRepo,
		newTestTokenService(t),
	)

	pair, err := service.Login(context.Background(), LoginInput{
		Login:    " user ",
		Password: password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)

	require.Equal(t, 42, createdToken.UserID)
	require.NotEmpty(t, createdToken.TokenHash)
	require.False(t, createdToken.ExpiresAt.IsZero())
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	credRepo := &credentialsRepositoryMock{
		getByLoginFunc: func(
			_ context.Context,
			_ string,
		) (domain.Credentials, error) {
			return domain.Credentials{}, core_errors.ErrNotFound
		},
	}

	service := NewAuthService(
		&tokenRepositoryMock{},
		credRepo,
		&userServiceMock{},
		newTestTokenService(t),
	)

	_, err := service.Login(context.Background(), LoginInput{
		Login:    "user",
		Password: "password123",
	})

	require.ErrorIs(t, err, core_errors.ErrInvalidCredentials)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte("correct-password"),
		bcrypt.MinCost,
	)
	require.NoError(t, err)

	credRepo := &credentialsRepositoryMock{
		getByLoginFunc: func(
			_ context.Context,
			_ string,
		) (domain.Credentials, error) {
			return domain.Credentials{
				UserID:       1,
				Login:        "user",
				PasswordHash: string(hash),
			}, nil
		},
	}

	service := NewAuthService(
		&tokenRepositoryMock{},
		credRepo,
		&userServiceMock{},
		newTestTokenService(t),
	)

	_, err = service.Login(context.Background(), LoginInput{
		Login:    "user",
		Password: "wrong-password",
	})

	require.ErrorIs(t, err, core_errors.ErrInvalidCredentials)
}

func TestAuthService_Refresh(t *testing.T) {
	oldToken := "old-refresh-token"

	oldTokenHash := HashRefreshToken(oldToken)

	tokenRepo := &tokenRepositoryMock{
		getRefreshTokenFunc: func(
			_ context.Context,
			tokenHash string,
		) (domain.RefreshToken, error) {
			require.Equal(t, oldTokenHash, tokenHash)

			return domain.RefreshToken{
				ID:        10,
				UserID:    42,
				TokenHash: oldTokenHash,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil
		},
	}

	var rotatedToken domain.RefreshToken
	tokenRepo.rotatedNewToken = &rotatedToken

	userRepo := &userServiceMock{
		getRoleFunc: func(
			_ context.Context,
			id int,
		) (string, error) {
			require.Equal(t, 42, id)
			return "buyer", nil
		},
	}

	service := NewAuthService(
		tokenRepo,
		&credentialsRepositoryMock{},
		userRepo,
		newTestTokenService(t),
	)

	pair, err := service.Refresh(
		context.Background(),
		oldToken,
	)

	require.NoError(t, err)

	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)

	require.Equal(t, 10, tokenRepo.rotatedOldTokenID)
	require.Equal(t, 42, rotatedToken.UserID)

	require.NotEqual(t, oldToken, pair.RefreshToken)
	require.Equal(t, HashRefreshToken(pair.RefreshToken), rotatedToken.TokenHash)
}

func TestAuthService_Refresh_ExpiredToken(t *testing.T) {
	tokenRepo := &tokenRepositoryMock{
		getRefreshTokenFunc: func(
			_ context.Context,
			_ string,
		) (domain.RefreshToken, error) {
			return domain.RefreshToken{
				ID:        1,
				UserID:    42,
				ExpiresAt: time.Now().Add(-time.Minute),
			}, nil
		},
	}

	service := NewAuthService(
		tokenRepo,
		&credentialsRepositoryMock{},
		&userServiceMock{},
		newTestTokenService(t),
	)

	_, err := service.Refresh(
		context.Background(),
		"refresh-token",
	)

	require.ErrorIs(t, err, core_errors.ErrInvalidRefreshToken)
	require.Zero(t, tokenRepo.rotatedOldTokenID)
}

func TestAuthService_Refresh_RevokedToken(t *testing.T) {
	revokedAt := time.Now().Add(-time.Minute)

	tokenRepo := &tokenRepositoryMock{
		getRefreshTokenFunc: func(
			_ context.Context,
			_ string,
		) (domain.RefreshToken, error) {
			return domain.RefreshToken{
				ID:        1,
				UserID:    42,
				ExpiresAt: time.Now().Add(time.Hour),
				RevokedAt: &revokedAt,
			}, nil
		},
	}

	service := NewAuthService(
		tokenRepo,
		&credentialsRepositoryMock{},
		&userServiceMock{},
		newTestTokenService(t),
	)

	_, err := service.Refresh(
		context.Background(),
		"refresh-token",
	)

	require.ErrorIs(t, err, core_errors.ErrInvalidRefreshToken)
	require.Zero(t, tokenRepo.rotatedOldTokenID)
}

func TestAuthService_Logout(t *testing.T) {
	token := "refresh-token"

	tokenRepo := &tokenRepositoryMock{
		getRefreshTokenFunc: func(
			_ context.Context,
			tokenHash string,
		) (domain.RefreshToken, error) {
			require.Equal(t, HashRefreshToken(token), tokenHash)

			return domain.RefreshToken{
				ID:        123,
				UserID:    42,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil
		},
	}

	service := NewAuthService(
		tokenRepo,
		&credentialsRepositoryMock{},
		&userServiceMock{},
		newTestTokenService(t),
	)

	err := service.Logout(
		context.Background(),
		token,
	)

	require.NoError(t, err)
	require.Equal(t, 123, tokenRepo.revokedTokenID)
}

func TestAuthService_Logout_NotFound(t *testing.T) {
	tokenRepo := &tokenRepositoryMock{
		getRefreshTokenFunc: func(
			_ context.Context,
			_ string,
		) (domain.RefreshToken, error) {
			return domain.RefreshToken{}, core_errors.ErrNotFound
		},
	}

	service := NewAuthService(
		tokenRepo,
		&credentialsRepositoryMock{},
		&userServiceMock{},
		newTestTokenService(t),
	)

	err := service.Logout(
		context.Background(),
		"invalid-token",
	)

	require.ErrorIs(t, err, core_errors.ErrInvalidRefreshToken)
}

func TestAuthService_Logout_RevokedToken(t *testing.T) {
	revokedAt := time.Now().Add(-time.Minute)

	tokenRepo := &tokenRepositoryMock{
		getRefreshTokenFunc: func(
			_ context.Context,
			_ string,
		) (domain.RefreshToken, error) {
			return domain.RefreshToken{
				ID:        123,
				UserID:    42,
				ExpiresAt: time.Now().Add(time.Hour),
				RevokedAt: &revokedAt,
			}, nil
		},
	}

	service := NewAuthService(
		tokenRepo,
		&credentialsRepositoryMock{},
		&userServiceMock{},
		newTestTokenService(t),
	)

	err := service.Logout(
		context.Background(),
		"refresh-token",
	)

	require.ErrorIs(t, err, core_errors.ErrInvalidRefreshToken)
	require.Zero(t, tokenRepo.revokedTokenID)
}
