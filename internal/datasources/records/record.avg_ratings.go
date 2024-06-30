package records

type AvgRatings struct {
	ID           uint    `db:"id"`
	ProductID    string  `db:"product_id"`
	Overall      float32 `db:"overall"`
	Quality      float32 `db:"quality"`
	Value        float32 `db:"value"`
	TotalReviews int     `db:"total_reviews"`
	Star5        int     `db:"star5"`
	Star4        int     `db:"star4"`
	Star3        int     `db:"star3"`
	Star2        int     `db:"star2"`
	Star1        int     `db:"star1"`
	TimeBase
}
