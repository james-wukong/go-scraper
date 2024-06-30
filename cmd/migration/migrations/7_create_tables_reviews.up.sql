BEGIN;

CREATE TABLE IF NOT EXISTS reviews(
    id SERIAL PRIMARY KEY,
    product_id uuid NOT NULL REFERENCES products (id),
    title VARCHAR(100) DEFAULT NULL,
    body TEXT DEFAULT NULL,
    author VARCHAR(100) DEFAULT NULL,
    author_id BIGINT DEFAULT NULL,
    is_recommended BOOLEAN DEFAULT false,
    rating SMALLINT NOT NULL DEFAULT 0,
    quality SMALLINT NOT NULL DEFAULT 0,
    value SMALLINT NOT NULL DEFAULT 0,
    count_helpful INTEGER NOT NULL DEFAULT 0,
    count_not_helpful INTEGER NOT NULL DEFAULT 0,
    fetch_ts VARCHAR(50) DEFAULT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    deleted_at timestamptz 
);

-- Add comments to the columns
COMMENT ON COLUMN reviews.id IS 'The primary key of reviews table';
COMMENT ON COLUMN reviews.product_id IS 'The foreign key references to products table';
COMMENT ON COLUMN reviews.rating IS 'The overall rating of the product';
COMMENT ON COLUMN reviews.quality IS 'The quality rating of the product';
COMMENT ON COLUMN reviews.value IS 'The value rating of the product';
COMMENT ON COLUMN reviews.title IS 'The title of the review';
COMMENT ON COLUMN reviews.body IS 'The body of the review';
COMMENT ON COLUMN reviews.author IS 'The author of the review';
COMMENT ON COLUMN reviews.author_id IS 'The author id like attribute of the review';
COMMENT ON COLUMN reviews.is_recommended IS 'is it recommended by author';
COMMENT ON COLUMN reviews.count_helpful IS 'people that think this review is helpful';
COMMENT ON COLUMN reviews.count_not_helpful IS 'people that think this review is not helpful';
COMMENT ON COLUMN reviews.fetch_ts IS 'based on created_at, how long has it been created';

COMMIT;