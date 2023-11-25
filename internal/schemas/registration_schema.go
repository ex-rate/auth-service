package schema

type Registration struct {
	Email          string `db:"email" json:"email"`
	PhoneNumber    string `db:"phone_number" json:"phone_number"`
	HashedPassword string `db:"hash_password" json:"hash_password"`
	Username       string `db:"username" json:"username"`
	FullName       string `db:"fullname" json:"fullname"`
}
