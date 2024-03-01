CREATE TYPE "Currency" AS ENUM (
  'USD',
  'EUR'
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar UNIQUE NOT NULL,
  "balance" bigint NOT NULL,
  "currency" "Currency" NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT current_timestamp
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT current_timestamp
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT current_timestamp
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD CONSTRAINT "fk_account_id" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD CONSTRAINT "fk_from_account_id" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD CONSTRAINT "fk_to_account_id" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
