package schema

// Для выдачи токена при регистрации / авторизации
type Token struct {
	RefreshToken string `json:"refresh-token"`
	AccessToken  string `json:"access-token"`
}

// Для запроса восстановления токена по адресу PUT /restore_token
type RestoreToken struct {
	RefreshToken string `json:"refresh-token"`
	AccessToken  string
}
