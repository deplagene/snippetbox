-- +goose Up
ALTER TABLE "users" RENAME COLUMN "password" TO "hashed_password";
ALTER TABLE "users" ALTER COLUMN "hashed_password" TYPE BYTEA USING ("hashed_password"::bytea);

ALTER TABLE "snippets" ADD COLUMN "user_id" uuid;
ALTER TABLE "snippets" ADD CONSTRAINT "snippets_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
