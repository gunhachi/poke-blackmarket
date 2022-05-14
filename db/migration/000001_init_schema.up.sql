CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL, 
  "user_role" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()) NOT NULL
);

CREATE TABLE "poke_orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "total_price" bigint NOT NULL,
  "order_detail" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()) NOT NULL
);

CREATE TABLE "poke_products" (
  "id" bigserial PRIMARY KEY,
  "poke_name" varchar NOT NULL,
  "status" varchar NOT NULL,
  "poke_price" bigint NOT NULL,
  "poke_stock" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now()) NOT NULL
);

CREATE INDEX ON "users" ("user_name");

CREATE INDEX ON "poke_orders" ("user_id");

CREATE INDEX ON "poke_orders" ("product_id");

CREATE INDEX ON "poke_products" ("poke_name");

COMMENT ON COLUMN "poke_orders"."quantity" IS 'must be positive';

COMMENT ON COLUMN "poke_orders"."total_price" IS 'must be positive';

COMMENT ON COLUMN "poke_products"."poke_price" IS 'must be positive';

COMMENT ON COLUMN "poke_products"."poke_stock" IS 'must be positive';

ALTER TABLE "poke_orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "poke_orders" ADD FOREIGN KEY ("product_id") REFERENCES "poke_products" ("id");
