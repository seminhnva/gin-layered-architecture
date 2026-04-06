CREATE TABLE
	IF NOT EXISTS users (
		user_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_name VARCHAR(50) UNIQUE NOT NULL ,
        name VARCHAR(100) NOT NULL,
		email VARCHAR(50) UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		deleted_at TIMESTAMPTZ DEFAULT NULL
	);


CREATE INDEX idx_users_created_at ON users (created_at);