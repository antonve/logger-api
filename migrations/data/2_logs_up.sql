CREATE SEQUENCE logs_seq;

CREATE TYPE activity AS ENUM ('FLASHCARDS','TEXTBOOK','READING','LISTENING','TRANSLATION','GRAMMAR','OTHER');
CREATE TYPE language AS ENUM ('JA','KR','ZH','DE');

CREATE TABLE logs (
  id bigint check (id > 0) NOT NULL DEFAULT NEXTVAL ('logs_seq'),
  user_id bigint NOT NULL,
  language language NOT NULL,
  date date NOT NULL,
  duration bigint check (duration > 0) NOT NULL,
  activity activity NOT NULL,
  notes jsonb,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY (id)
);

ALTER SEQUENCE logs_seq RESTART WITH 1;
