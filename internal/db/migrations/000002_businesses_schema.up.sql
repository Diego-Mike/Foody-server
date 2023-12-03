CREATE TABLE "businesses" (
    "business_id" bigserial PRIMARY KEY,
    "name" varchar(30) NOT NULL,
    "city" varchar(50) NOT NULL,
    "address" varchar(100) NOT NULL,
    "latitude" varchar(100) NOT NULL,
    "longitude" varchar(100) NOT NULL,
    "presentation" varchar(300) NOT NULL, 
    "clients_max_amount" smallint,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "business_members" (
    "business_id" bigint REFERENCES "businesses" ("business_id"),
    "user_id" bigint REFERENCES "users" ("user_id"),
    "business_position" varchar(50) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    PRIMARY KEY ("business_id", "user_id")
);


COMMENT ON COLUMN "businesses"."presentation" IS  'max amount of people that can order food in this business';
COMMENT ON COLUMN "businesses"."clients_max_amount" IS 'number of people that can be inside the business, field is not mandatory';
COMMENT ON TABLE "business_members" IS 'primary key 🔑 is composed by business_id and user_id';
COMMENT ON COLUMN "business_members"."business_position" IS 'posible positions for the employee, dueño, administrador, empleado';
