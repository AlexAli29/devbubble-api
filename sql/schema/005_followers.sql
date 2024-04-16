-- +goose Up
CREATE TABLE "user_follows" (
    "followerId" UUID NOT NULL,
    "followeeId" UUID NOT NULL,
    PRIMARY KEY ("followerId", "followeeId"),
    FOREIGN KEY ("followerId") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY ("followeeId") REFERENCES "users"("id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "user_follows";