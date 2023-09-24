CREATE TABLE "users" (
  "user_id" bigserial UNIQUE,
  "social_id" varchar UNIQUE,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "picture" varchar NOT NULL,
  "provider" varchar NOT NULL,
  "registered_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY("user_id", "social_id")
);

COMMENT ON TABLE "users" IS 'primary key is composed by user_id and social_id, both of those properties are unique';
COMMENT ON COLUMN users.social_id IS 'this field is for the unique id provided by google, or facebook, or tik tok etc etc.. to identify a user';

CREATE TABLE "sessions" (
  "user_id_session" bigint REFERENCES users (user_id),
  "valid" boolean NOT NULL,
  "user_agent" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY("user_id_session")
);
