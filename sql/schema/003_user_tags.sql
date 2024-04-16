-- +goose Up
CREATE TABLE "user_tags" (
   "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
   "text" TEXT NOT NULL,
   "icon" TEXT NOT NULL
);

INSERT INTO "user_tags" ("text", "icon") VALUES 
('C#', 'c#'),
('JS', 'js'),
('Go', 'go'),
('Devops', 'devops'),
('Kotlin', 'kotlin'),
('TS', 'ts'),
('C++', 'c++'),
('Python', 'python');

CREATE TABLE "user_user_tags" (
   "userId" uuid NOT NULL,
   "userTagId" uuid NOT NULL,
   PRIMARY KEY ("userId", "userTagId"),
   FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE,
   FOREIGN KEY ("userTagId") REFERENCES "user_tags"("id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE "user_user_tags";
DROP TABLE "user_tags";