-- Create the initial schema

-- +goose Up

-- auth
CREATE TABLE "users" (
  "id" SERIAL NOT NULL,
  "email" VARCHAR(256) NOT NULL,
  "name" VARCHAR(128) NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "unique_user_emails" ON "users" ("email") WHERE ("deleted_at" IS NULL AND "email" != '');

CREATE TABLE "sessions" (
  "key" VARCHAR NOT NULL,
  "user_id" INTEGER NOT NULL REFERENCES users("id") ON DELETE CASCADE,
  "expires_at" TIMESTAMP NOT NULL,
  PRIMARY KEY ("key")
);

-- things
CREATE TABLE "things" (
  "id" SERIAL NOT NULL,
  "content" JSON NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP,
  PRIMARY KEY ("id")
);

-- +goose Down
DROP TABLE IF EXISTS "things";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "users";
