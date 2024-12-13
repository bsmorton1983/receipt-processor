CREATE TABLE "receipts" (
  "id" bigserial PRIMARY KEY,
  "retailer" varchar NOT NULL,
  "purchase_date" varchar NOT NULL,
  "purchase_time" varchar NOT NULL,
  "creation_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "receipt_items" (
  "id" bigserial PRIMARY KEY,
  "receipt_id" bigint NOT NULL,
  "short_description" varchar NOT NULL,
  "price" double precision NOT NULL,
  "creation_time" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "receipts" ("retailer");

CREATE INDEX ON "receipt_items" ("receipt_id");

ALTER TABLE "receipt_items" ADD FOREIGN KEY ("receipt_id") REFERENCES "receipts" ("id");