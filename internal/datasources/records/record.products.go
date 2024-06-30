package records

type Products struct {
	Id         string      `db:"id"`
	CategoryId int         `db:"category_id"`
	Name       string      `db:"name"`
	Sku        string      `db:"sku"`
	Model      string      `db:"model"`
	ProdId     string      `db:"prod_id"`
	Price      float32     `db:"price"`
	Source     int         `db:"source"`
	UrlLink    string      `db:"url_link"`
	ImageLink  string      `db:"img_link"`
	Detail     *DetailBase `db:"detail" json:"detail"`
	Spec       *SpecBase   `db:"specification" json:"specification"`
	TimeBase
}
