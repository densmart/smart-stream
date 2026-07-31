-- Write your migrate up statements here

-- Change size column type from int to bigint to support large media files (>2GB)
ALTER TABLE "media" ALTER COLUMN "size" TYPE bigint;

---- create above / drop below ----

-- Revert size column type back to int
ALTER TABLE "media" ALTER COLUMN "size" TYPE int;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.