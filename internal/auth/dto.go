package auth

type UserRegister struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	DisplayName  string `json:"display_name"`
	PasswordHash string `json:"-"`
	PasswordSalt string `json:"-"`
}

type UserLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	ID           uint64 `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	PasswordSalt string `json:"-"`
	DisplayName  string `json:"display_name"`
}
