BEGIN;

CREATE TABLE IF NOT EXISTS avg_ratings(
    id SERIAL PRIMARY KEY,
    product_id uuid NOT NULL REFERENCES products (id),
    overall NUMERIC(2,1) NOT NULL DEFAULT 0,
    quality NUMERIC(2,1) NOT NULL DEFAULT 0,
    value NUMERIC(2,1) NOT NULL DEFAULT 0,
    total_reviews INT DEFAULT 0,
    star5 INTEGER NOT NULL DEFAULT 0,
    star4 INTEGER NOT NULL DEFAULT 0,
    star3 INTEGER NOT NULL DEFAULT 0,
    star2 INTEGER NOT NULL DEFAULT 0,
    star1 INTEGER NOT NULL DEFAULT 0,
    
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    deleted_at timestamptz 
);

-- Add comments to the columns
COMMENT ON COLUMN avg_ratings.id IS 'The primary key of discount_histories table';
COMMENT ON COLUMN avg_ratings.product_id IS 'The foreign key references to products table';
COMMENT ON COLUMN avg_ratings.overall IS 'The overall rating of the product';
COMMENT ON COLUMN avg_ratings.quality IS 'The quality rating of the product';
COMMENT ON COLUMN avg_ratings.value IS 'The value rating of the product';

COMMIT;