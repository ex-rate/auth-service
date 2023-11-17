package entities

type User struct {
	ID          string `db:"user_id"`
	Username    string `db:"username"`
	Password    string `db:"hash_password"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	FullName    string `db:"fullname"`
}
