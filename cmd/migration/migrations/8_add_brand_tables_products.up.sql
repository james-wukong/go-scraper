BEGIN;

ALTER TABLE products 
ADD COLUMN brand VARCHAR(100) DEFAULT NULL;

COMMIT;