CREATE TABLE "business_schedule" (
    "business_id" bigserial NOT NULL REFERENCES businesses (business_id),
    "day_of_week" smallint NOT NULL,
    "opening_hour" time NOT NULL,
    "closing_hour" time NOT NULL
);

COMMENT ON COLUMN "business_schedule"."day_of_week" IS '1 = lunes, 2 = martes';
CREATE INDEX idx_schedule_business_id ON business_schedule(business_id);

CREATE TABLE "business_food" (
    "business_id" bigint NOT NULL REFERENCES businesses (business_id),
    "food_id" bigserial,
    "food_img" varchar NOT NULL,
    "food_title" varchar(55) NOT NULL,
    "food_description" varchar(150),
    "food_price" bigint NOT NULL,
    "food_available_per_day" smallint,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY("business_id", "food_id")
);

COMMENT ON COLUMN "business_food"."food_available_per_day" IS 'number of pieces of this food that can be sold per day, it is not mandatory';
