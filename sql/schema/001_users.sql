-- +goose Up
CREATE TABLE "users" (
    "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,   
    "email" TEXT NOT NULL UNIQUE,
    "description" TEXT,
    "name" TEXT NOT NULL,
    "isVerified" BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE "users";