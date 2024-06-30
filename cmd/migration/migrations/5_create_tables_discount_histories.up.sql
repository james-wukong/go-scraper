BEGIN;

CREATE TABLE IF NOT EXISTS discount_histories(
    id SERIAL PRIMARY KEY,
    product_id uuid NOT NULL REFERENCES products (id),
    price NUMERIC DEFAULT NULL,
    save_amount NUMERIC DEFAULT NULL,
    save_percent NUMERIC(4, 2) DEFAULT NULL,
    duration VARCHAR(50) DEFAULT NULL,
    started_at timestamptz NOT NULL,
    ended_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    deleted_at timestamptz 
);

-- Add comments to the columns
COMMENT ON COLUMN discount_histories.id IS 'The primary key of discount_histories table';
COMMENT ON COLUMN discount_histories.product_id IS 'The foreign key references to products table';
COMMENT ON COLUMN discount_histories.price IS 'The price of the product, without discount and is the value in accordance with fetch time';
COMMENT ON COLUMN discount_histories.save_amount IS 'The save_amount of the product in this promotion';
COMMENT ON COLUMN discount_histories.save_percent IS 'The save_percent of the product in this promotion';
COMMENT ON COLUMN discount_histories.duration IS 'The duration of promotion';
COMMENT ON COLUMN discount_histories.started_at IS 'promotion started at';
COMMENT ON COLUMN discount_histories.ended_at IS 'promotion ended at';

COMMIT;