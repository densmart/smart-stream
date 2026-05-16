-- Write your migrate up statements here

-- Add parent_id column to playlists table for hierarchy support
ALTER TABLE "playlists" ADD COLUMN "parent_id" UUID REFERENCES "playlists" ("id") ON DELETE CASCADE DEFAULT NULL;

-- Add has_children column to track if playlist has child playlists
ALTER TABLE "playlists" ADD COLUMN "has_children" BOOLEAN DEFAULT FALSE NOT NULL;

-- Create indexes for better query performance
CREATE INDEX playlists_parent_id ON "playlists" ("parent_id");
CREATE INDEX playlists_has_children ON "playlists" ("has_children");

---- create above / drop below ----

-- Drop indexes
DROP INDEX playlists_has_children;
DROP INDEX playlists_parent_id;

-- Drop columns
ALTER TABLE "playlists" DROP COLUMN "has_children";
ALTER TABLE "playlists" DROP COLUMN "parent_id";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
