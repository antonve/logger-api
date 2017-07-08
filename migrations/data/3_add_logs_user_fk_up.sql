ALTER TABLE logs
ADD CONSTRAINT logs_user_fk FOREIGN KEY (user_id) REFERENCES users (id);
