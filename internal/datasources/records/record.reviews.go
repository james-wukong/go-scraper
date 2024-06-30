package records

type Reviews struct {
	Id              uint   `db:"id"`
	ProductId       string `db:"product_id"`
	Title           string `db:"title"`
	Body            string `db:"body"`
	Author          string `db:"author"`
	AuthorId        uint   `db:"author_id"`
	IsRecommended   bool   `db:"is_recommended"`
	Rating          uint8  `db:"rating"`
	Quality         uint8  `db:"quality"`
	Value           uint8  `db:"value"`
	CountHelpful    uint   `db:"count_helpful"`
	CountNotHelpful uint   `db:"count_not_helpful"`
	FetchTs         string `db:"fetch_ts"`
	TimeBase
}
