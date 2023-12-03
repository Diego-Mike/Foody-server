CREATE TABLE "business_reservations" (
  "reservation_id" bigserial PRIMARY KEY,
  "business_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "order_schedule" timestamptz,
  "accepted" boolean NOT NULL DEFAULT (false),
  "created_at" timestamptz DEFAULT (now())
);

-- id for this table ?¿
CREATE TABLE "business_reservations_state" (
  "reservation_id" bigint NOT NULL,
  "cancelled_by_client" boolean NOT NULL,
  "cancelled_by_business" boolean NOT NULL,
  "reason_for_cancellation" varchar(300) NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

-- id for this table ?¿
CREATE TABLE "reserve_food" (
  "reservation_id" bigint NOT NULL,
  "food_id" bigint NOT NULL,
  "amount" smallint NOT NULL,
  "details" varchar(250),
  "created_at" timestamptz DEFAULT (now()),
  PRIMARY KEY ("reservation_id", "food_id")
);

-- id for this table ?¿
CREATE TABLE "business_reservations_notificacions" (
  "reservation_id" bigint NOT NULL,
  "notification_title" varchar NOT NULL,
  "notification_description" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "business_reservations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
ALTER TABLE "business_reservations" ADD FOREIGN KEY ("business_id") REFERENCES "businesses" ("business_id");

ALTER TABLE "business_reservations_state" ADD FOREIGN KEY ("reservation_id") REFERENCES "business_reservations" ("reservation_id");

ALTER TABLE "reserve_food" ADD FOREIGN KEY ("reservation_id") REFERENCES "business_reservations" ("reservation_id");
ALTER TABLE "reserve_food" ADD FOREIGN KEY ("food_id") REFERENCES "business_food" ("food_id");

ALTER TABLE "business_reservations_notificacions" ADD FOREIGN KEY ("reservation_id") REFERENCES "business_reservations" ("reservation_id");
