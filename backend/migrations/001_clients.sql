-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create clients table
CREATE TABLE "clients" (
    "id"                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at"        timestamptz,
    "updated_at"        timestamptz,
    "login"             varchar(255) UNIQUE NOT NULL,
    "password"          varchar(255) NOT NULL,
    "email"             varchar(255) UNIQUE NULL,
    "current_version"   varchar(4096) NULL,
    "is_active"         boolean DEFAULT false,
    "last_login_at"     timestamptz DEFAULT NULL
);
CREATE INDEX clients_created_at ON "clients" ("created_at");
CREATE INDEX clients_login ON "clients" ("login");
CREATE INDEX clients_email ON "clients" ("email");
CREATE INDEX clients_is_active ON "clients" ("is_active");
CREATE INDEX clients_last_login_at ON "clients" ("last_login_at");

---- create above / drop below ----

DROP TABLE "clients";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
