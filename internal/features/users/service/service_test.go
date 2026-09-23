package service

import (
	"context"
	"errors"
	"testing"

	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	users  map[int]domain.User
	nextID int

	createErr error
	getErr    error
	getAllErr error
	patchErr  error
	deleteErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:  make(map[int]domain.User),
		nextID: 1,
	}
}

func (r *mockUserRepo) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if r.createErr != nil {
		return domain.User{}, r.createErr
	}

	user.ID = r.nextID
	user.Version = 1

	r.users[user.ID] = user
	r.nextID++

	return user, nil
}

func (r *mockUserRepo) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	if r.getAllErr != nil {
		return nil, r.getAllErr
	}

	users := make([]domain.User, 0, len(r.users))

	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *mockUserRepo) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	if r.getErr != nil {
		return domain.User{}, r.getErr
	}

	user, ok := r.users[id]
	if !ok {
		return domain.User{}, core_errors.ErrNotFound
	}

	return user, nil
}

func (r *mockUserRepo) DeleteUser(
	ctx context.Context,
	id int,
) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}

	if _, ok := r.users[id]; !ok {
		return core_errors.ErrNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *mockUserRepo) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	if r.patchErr != nil {
		return domain.User{}, r.patchErr
	}

	user, ok := r.users[id]
	if !ok {
		return domain.User{}, core_errors.ErrNotFound
	}

	if patch.Name.Set {
		user.Name = *patch.Name.Value
	}

	if patch.Phone.Set {
		user.Phone = patch.Phone.Value
	}

	if patch.Role.Set {
		user.Role = domain.UserRole(*patch.Role.Value)
	}

	user.Version++
	r.users[id] = user

	return user, nil
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "79991234567"

		user, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name:  "Nikita",
				Phone: &phone,
				Role:  domain.UserBuyer,
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, int64(1), user.Version)
		assert.Equal(t, "Nikita", user.Name)
		assert.Equal(t, &phone, user.Phone)
		assert.Equal(t, domain.UserBuyer, user.Role)
	})

	t.Run("blank name", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name: "   ",
				Role: domain.UserBuyer,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidName)
	})

	t.Run("empty role", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name: "Nikita",
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidRole)
	})

	t.Run("nil phone", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		user, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name:  "Nikita",
				Phone: nil,
				Role:  domain.UserBuyer,
			},
		)

		require.NoError(t, err)
		assert.Nil(t, user.Phone)
	})

	t.Run("phone too short", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "7999123456"

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name:  "Nikita",
				Phone: &phone,
				Role:  domain.UserBuyer,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("phone too long", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "799912345678"

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name:  "Nikita",
				Phone: &phone,
				Role:  domain.UserBuyer,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("phone contains letters", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "7999123456a"

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name:  "Nikita",
				Phone: &phone,
				Role:  domain.UserBuyer,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("phone contains spaces", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "7999 123456"

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name:  "Nikita",
				Phone: &phone,
				Role:  domain.UserBuyer,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.createErr = errors.New("database error")

		service := NewUserService(repo)

		_, err := service.CreateUser(
			context.Background(),
			domain.User{
				Name: "Nikita",
				Role: domain.UserBuyer,
			},
		)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.users[1] = domain.User{
			ID:   1,
			Name: "Nikita",
			Role: domain.UserBuyer,
		}

		service := NewUserService(repo)

		user, err := service.GetUser(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "Nikita", user.Name)
	})

	t.Run("zero id", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.GetUser(context.Background(), 0)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("negative id", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.GetUser(context.Background(), -1)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.GetUser(context.Background(), 1)

		require.ErrorIs(t, err, core_errors.ErrNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.getErr = errors.New("database error")

		service := NewUserService(repo)

		_, err := service.GetUser(context.Background(), 1)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestUserService_GetUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepo()

		repo.users[1] = domain.User{
			ID:   1,
			Name: "User 1",
		}
		repo.users[2] = domain.User{
			ID:   2,
			Name: "User 2",
		}

		service := NewUserService(repo)

		users, err := service.GetUsers(
			context.Background(),
			nil,
			nil,
		)

		require.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("negative limit", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		limit := -1

		_, err := service.GetUsers(
			context.Background(),
			&limit,
			nil,
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidQueryParam)
	})

	t.Run("negative offset", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		offset := -1

		_, err := service.GetUsers(
			context.Background(),
			nil,
			&offset,
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidQueryParam)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.getAllErr = errors.New("database error")

		service := NewUserService(repo)

		_, err := service.GetUsers(
			context.Background(),
			nil,
			nil,
		)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestUserService_PatchUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepo()

		repo.users[1] = domain.User{
			ID:      1,
			Version: 1,
			Name:    "Old",
			Role:    domain.UserBuyer,
		}

		service := NewUserService(repo)

		name := "New"

		user, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: &name,
				},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, "New", user.Name)
		assert.Equal(t, int64(2), user.Version)
	})

	t.Run("patch phone", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.users[1] = domain.User{
			ID:   1,
			Name: "Nikita",
		}

		service := NewUserService(repo)

		phone := "79991234567"

		user, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Phone: domain.Nullable[string]{
					Set:   true,
					Value: &phone,
				},
			},
		)

		require.NoError(t, err)
		require.NotNil(t, user.Phone)
		assert.Equal(t, phone, *user.Phone)
	})

	t.Run("patch phone to null", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.users[1] = domain.User{
			ID:   1,
			Name: "Nikita",
		}

		service := NewUserService(repo)

		user, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Phone: domain.Nullable[string]{
					Set:   true,
					Value: nil,
				},
			},
		)

		require.NoError(t, err)
		assert.Nil(t, user.Phone)
	})

	t.Run("patch role", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.users[1] = domain.User{
			ID:   1,
			Name: "Nikita",
			Role: domain.UserBuyer,
		}

		service := NewUserService(repo)

		role := string(domain.UserSeller)

		user, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Role: domain.Nullable[string]{
					Set:   true,
					Value: &role,
				},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, domain.UserSeller, user.Role)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.PatchUser(
			context.Background(),
			0,
			domain.UserPatch{},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("name null", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: nil,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrNullNotAllowed)
	})

	t.Run("blank name", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		name := "   "

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: &name,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidName)
	})

	t.Run("phone too short", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "7999123456"

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Phone: domain.Nullable[string]{
					Set:   true,
					Value: &phone,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("phone too long", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "799912345678"

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Phone: domain.Nullable[string]{
					Set:   true,
					Value: &phone,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("phone contains non digits", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		phone := "7999123456a"

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Phone: domain.Nullable[string]{
					Set:   true,
					Value: &phone,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPhone)
	})

	t.Run("role null", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Role: domain.Nullable[string]{
					Set:   true,
					Value: nil,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrNullNotAllowed)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.patchErr = errors.New("database error")

		service := NewUserService(repo)

		name := "New"

		_, err := service.PatchUser(
			context.Background(),
			1,
			domain.UserPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: &name,
				},
			},
		)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.users[1] = domain.User{ID: 1}

		service := NewUserService(repo)

		err := service.DeleteUser(context.Background(), 1)

		require.NoError(t, err)
		assert.NotContains(t, repo.users, 1)
	})

	t.Run("zero id", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		err := service.DeleteUser(context.Background(), 0)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("negative id", func(t *testing.T) {
		repo := newMockUserRepo()
		service := NewUserService(repo)

		err := service.DeleteUser(context.Background(), -1)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.deleteErr = errors.New("database error")

		service := NewUserService(repo)

		err := service.DeleteUser(context.Background(), 1)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}
