CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    username varchar(255) UNIQUE NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);
