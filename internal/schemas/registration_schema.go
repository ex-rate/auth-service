package schema

type Registration struct {
	Email          string `db:"email" json:"email"`
	PhoneNumber    string `db:"phone_number" json:"phone_number"`
	HashedPassword string `db:"hash_password" json:"hash_password"`
	Username       string `db:"username" json:"username"`
	LastName       string `db:"last_name" json:"last_name"`
	FirstName      string `db:"first_name" json:"first_name"`
	Patronymic     string `db:"patronymic" json:"patronymic"`
}
