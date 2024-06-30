BEGIN;

CREATE TABLE IF NOT EXISTS categories(
    id SERIAL PRIMARY KEY,
    parent_id INTEGER DEFAULT NULL REFERENCES categories,
    name VARCHAR(50) NOT NULL,
    level SMALLINT DEFAULT 0,
    url VARCHAR(250) DEFAULT NULL,
    platform SMALLINT DEFAULT 0
);

CREATE INDEX  IF NOT EXISTS idx_categories_name ON categories (name);

-- Add comments to the columns
COMMENT ON COLUMN categories.id IS 'The primary key of categories table';
COMMENT ON COLUMN categories.parent_id IS 'The parent id of the category, default null';
COMMENT ON COLUMN categories.name IS 'The name of the category';
COMMENT ON COLUMN categories.level IS 'The level of the category, starts from 0';
COMMENT ON COLUMN categories.url IS 'The url of the category';
COMMENT ON COLUMN categories.platform IS 'The platform of the category, starts from 100, check constants for details';

COMMIT;