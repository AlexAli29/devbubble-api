-- +goose Up
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,   
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE users;
