-- Write your migrate up statements here

-- create playlist table
CREATE TABLE "playlists" (
    "id"            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at"    timestamptz,
    "updated_at"    timestamptz,
    "name"          varchar(255) NOT NULL,
    "type"          varchar(255) NOT NULL, -- e.g. "franchise", "series", "favourites", etc.
    "poster"        varchar(4096) NULL -- path to the poster image on disk
);

CREATE INDEX playlists_created_at ON "playlists" ("created_at");
CREATE INDEX playlists_name ON "playlists" ("name");
CREATE INDEX playlists_type ON "playlists" ("type");

-- create media table
CREATE TABLE "media" (
    "id"            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at"    timestamptz,
    "updated_at"    timestamptz,
    "playlist_id"   UUID REFERENCES playlists ("id") ON DELETE SET NULL DEFAULT NULL,
    "name"          varchar(255) NOT NULL,
    "poster"        varchar(4096) NULL, -- path to the poster image on disk
    "format"        varchar(255) NOT NULL, -- e.g. "video/mp4", "video/mkv", etc.
    "path"          varchar(4096) NOT NULL, -- path to the media file on disk
    "duration"      int NOT NULL, -- in seconds
    "size"          int NOT NULL, -- in bytes
    "order"         int DEFAULT 0
);

CREATE INDEX media_created_at ON "media" ("created_at");
CREATE INDEX media_name ON "media" ("name");
CREATE INDEX media_playlist_id ON "media" ("playlist_id");
CREATE INDEX media_format ON "media" ("format");
CREATE INDEX media_order ON "media" ("order");

---- create above / drop below ----

DROP TABLE "media";
DROP TABLE "playlists";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
