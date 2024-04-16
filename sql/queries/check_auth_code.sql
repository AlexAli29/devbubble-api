-- name: CheckAuthCode :one
WITH matched_user AS (
    SELECT "auth_codes"."userId"
    FROM "auth_codes"
    JOIN "users" ON "auth_codes"."userId" = "users"."id"
    WHERE "users"."email" = $1
      AND "users"."isVerified"=true
      AND "auth_codes"."code" = $2
      AND "auth_codes"."updatedAt" > (CURRENT_TIMESTAMP - INTERVAL '10 MINUTES')
    LIMIT 1
),
updated AS (
    UPDATE "auth_codes"
    SET "code" = floor(random() * 900000 + 100000) -- Generates a new code        
    WHERE "userId" IN (SELECT "userId" FROM matched_user)
    RETURNING "userId"
)
SELECT "userId" FROM updated;