CREATE SEQUENCE refresh_tokens_seq;

CREATE TABLE refresh_tokens (
  id bigint check (id > 0) NOT NULL DEFAULT NEXTVAL ('refresh_tokens_seq'),
  user_id bigint NOT NULL REFERENCES users (id),
  device_id varchar(32) NOT NULL,
  refresh_token bytea NOT NULL,
  created_at timestamp NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  updated_at timestamp NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  invalidated_at timestamp DEFAULT NULL,
  PRIMARY KEY (id)
);

ALTER SEQUENCE refresh_tokens_seq RESTART WITH 1;
