package model

import (
	"fmt"
	"strconv"
	"time"
)

// User  用户表
type User struct {
	Id        *int       `json:"id" db:"id,pk" uri:"id"` // 编号
	Username  *string    `json:"username" db:"username"` // 用户账户
	Password  *string    `json:"password" db:"password"` // 用户密码
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

func (*User) TableName() string {
	return "user"
}

func (o *User) CacheKey() string {
	return o.TableName() + ":id:" + strconv.Itoa(*o.Id)
}
func UserCacheKey(id int) string {
	return fmt.Sprintf("user:id:%d", id)
}
