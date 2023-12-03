CREATE TABLE "business_food" (
    "business_id" bigint NOT NULL REFERENCES businesses (business_id),
    "food_id" bigserial PRIMARY KEY,
    "food_img" varchar NOT NULL,
    "food_title" varchar(55) NOT NULL,
    "food_description" varchar(150),
    "food_price" bigint NOT NULL,
    "food_available_per_day" smallint,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "business_food"."food_available_per_day" IS 'number of pieces of this food that can be sold per day, it is not mandatory';
