BEGIN TRANSACTION ;

	
	CREATE TABLE IF NOT EXISTS users
	(
		id bigint NOT NULL,
		username VARCHAR(255) NOT NULL,
		password bytea NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	);
	
	ALTER TABLE users ADD CONSTRAINT username_unq UNIQUE(username);

	CREATE INDEX users_username_key ON users USING HASH (username);

COMMIT ;