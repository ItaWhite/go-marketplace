package repository

import "time"

type UserModel struct {
	ID        int
	Version   int64
	Name      string
	Phone     *string
	Role      string
	CreatedAt time.Time
}
