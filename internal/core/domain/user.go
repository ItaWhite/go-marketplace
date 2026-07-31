package domain

import "time"

type UserRole string

const (
	UserBuyer  UserRole = "buyer"
	UserSeller UserRole = "seller"
)

type User struct {
	ID        int       `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (r UserRole) IsValid() bool {
	switch r {
	case UserBuyer, UserSeller:
		return true
	default:
		return false
	}
}

type UserPatch struct {
	Name  Nullable[string]
	Phone Nullable[string]
	Role  Nullable[string]
}
