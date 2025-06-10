package input

type UserLogin struct {
	Username *string `json:"username" db:"username"`
	Password *string `json:"password" db:"password"`
}

type UserRegister struct {
	Username *string `json:"username" db:"username"`
	Password *string `json:"password" db:"password"`
}

type UserFilter struct {
	Username *string `json:"username" db:"username"`
}
