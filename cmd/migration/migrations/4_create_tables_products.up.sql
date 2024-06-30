BEGIN;

CREATE TABLE IF NOT EXISTS products(
    id uuid PRIMARY KEY,
    category_id INTEGER DEFAULT NULL REFERENCES categories (id),
    name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) DEFAULT NULL,
    model VARCHAR(50) DEFAULT NULL,
    prod_id VARCHAR(50) DEFAULT NULL,
    price NUMERIC DEFAULT NULL,
    source SMALLINT DEFAULT 0,
    url_link VARCHAR(255) DEFAULT NULL,
    img_link VARCHAR(255) DEFAULT NULL,
    detail JSONB DEFAULT NULL,
    specification JSONB DEFAULT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    deleted_at timestamptz 
);

CREATE INDEX  IF NOT EXISTS idx_products_name ON products (name);
CREATE INDEX  IF NOT EXISTS idx_products_sku ON products (sku);
CREATE INDEX  IF NOT EXISTS idx_products_source ON products (source);

-- Add comments to the columns
COMMENT ON COLUMN products.id IS 'The primary key of products table';
COMMENT ON COLUMN products.category_id IS 'The foreign key references to categories table';
COMMENT ON COLUMN products.name IS 'The name of the product';
COMMENT ON COLUMN products.sku IS 'The sku of the product';
COMMENT ON COLUMN products.model IS 'The model of the product';
COMMENT ON COLUMN products.price IS 'The price of the product, without discount and is the value in accordance with fetch time';
COMMENT ON COLUMN products.source IS 'The source platform of the product';
COMMENT ON COLUMN products.url_link IS 'The url of the product';
COMMENT ON COLUMN products.img_link IS 'The image url of the product';
COMMENT ON COLUMN products.detail IS 'The detail of the product, is a jsonb format: key => list';
COMMENT ON COLUMN products.specification IS 'The specification of the product, is a jsonb format: key => value';

COMMIT;