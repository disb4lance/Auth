package service

// DTO для входа
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// возвращаем клиенту только безопасные поля
type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// DTO для ответа после аутентификации
type AuthenticatedUser struct {
	User  UserDTO   `json:"user"`
	Token TokenPair `json:"token"`
}

type AuthService interface {
	Register(email, password string) (*UserDTO, error)
	Authenticate(creds Credentials) (*TokenResponse, error)
	Refresh(refreshToken string) (*TokenResponse, error)
}
