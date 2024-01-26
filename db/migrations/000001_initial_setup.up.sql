CREATE TABLE "users" (
  "id" UUID PRIMARY KEY,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "email" TEXT UNIQUE NOT NULL,
  "password" TEXT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE "tasks" (
  "id" UUID PRIMARY KEY,
  "title" TEXT NOT NULL,
  "due_date" TIMESTAMPTZ NOT NULL,
  "reminder_date" TIMESTAMPTZ,
  "description" TEXT,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

ALTER TABLE "tasks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
