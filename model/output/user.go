package output

import "time"

type UserOutput struct {
	Id        *int       `json:"id" db:"id" uri:"id"`
	Username  *string    `json:"username" db:"username"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}
