package records

type Categories struct {
	Id       int    `db:"id"`
	ParentId *int   `db:"parent_id"`
	Name     string `db:"name"`
	Level    uint   `db:"level"`
	Url      string `db:"url"`
	Platform uint   `db:"platform"`
}
