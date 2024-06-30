package records

type Users struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Active   bool   `db:"active"`
	RoleId   int    `db:"role_id"`
	TimeBase
}
