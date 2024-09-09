CREATE TABLE IF NOT EXISTS tokens (
    id serial PRIMARY KEY,
    token varchar NOT NULL,
    user_id serial NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE
);