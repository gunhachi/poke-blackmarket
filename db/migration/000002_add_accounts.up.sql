CREATE TABLE "accounts" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()) NOT NULL,
  "password_changet_at" timestamptz DEFAULT '0001-01-01 00:00:00+00:00'
);

ALTER TABLE "users" ADD FOREIGN KEY ("user_name") REFERENCES "accounts" ("username");

-- CREATE UNIQUE INDEX ON "users" ("user_name", "user_role");
ALTER TABLE "users" ADD CONSTRAINT "name_role_key" UNIQUE ("user_name","user_role");

