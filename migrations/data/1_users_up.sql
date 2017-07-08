CREATE SEQUENCE user_seq;

CREATE TYPE role AS ENUM ('ADMIN','USER','DISABLED');

CREATE TABLE users (
  id bigint check (id > 0) NOT NULL DEFAULT NEXTVAL ('user_seq'),
  username varchar(255) NOT NULL UNIQUE,
  display_name varchar(255) NOT NULL,
  password bytea NOT NULL,
  role role NOT NULL DEFAULT 'DISABLED',
  PRIMARY KEY (id)
);

ALTER SEQUENCE user_seq RESTART WITH 1;
