package schema

type Registration struct {
	Email          string `db:"email"`
	PhoneNumber    string `db:"phone_number"`
	HashedPassword string `db:"hash_password"`
	Username       string `db:"username"`
	FullName       string `db:"fullname"`
}
