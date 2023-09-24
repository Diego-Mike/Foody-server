CREATE TABLE "businesses" (
    "business_id" bigserial PRIMARY KEY,
    "name" varchar(30) NOT NULL,
    "city" varchar(50) NOT NULL,
    "address" varchar(100) NOT NULL,
    "latitude" varchar(100) NOT NULL,
    "longitude" varchar(100) NOT NULL,
    "ubication_photo" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "business_members" (
    "business_id" bigint,
    "user_id" bigint,
    "business_position" varchar(50) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    PRIMARY KEY ("business_id", "user_id")
);

ALTER TABLE "business_members" ADD FOREIGN KEY ("business_id") REFERENCES "businesses" ("business_id");
ALTER TABLE "business_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

COMMENT ON TABLE "business_members" IS 'primary key 🔑 is composed by business_id and user_id';
COMMENT ON COLUMN "business_members"."business_position" IS 'posible positions for the employee, dueño, administrador, empleado';
