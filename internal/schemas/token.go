package schema

type Token struct {
	RefreshToken string `json:"refresh-token"`
	AccessToken  string `json:"access-token"`
}
