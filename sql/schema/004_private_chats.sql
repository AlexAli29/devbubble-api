-- +goose Up
CREATE TABLE "private_chats" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY
);

CREATE TABLE "chat_participants" (
    "chatId" UUID NOT NULL,
    "userId" UUID NOT NULL,
    PRIMARY KEY ("chatId", "userId"),
    FOREIGN KEY ("chatId") REFERENCES "private_chats"("id") ON DELETE CASCADE,
    FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE TABLE "messages" (
    "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "text" TEXT NOT NULL,
    "userId" UUID NOT NULL,
    "chatId" UUID NOT NULL,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY ("chatId") REFERENCES "private_chats"("id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "messages";
DROP TABLE IF EXISTS "chat_participants";
DROP TABLE IF EXISTS "private_chats";