-- name: VerifyUser :one
WITH matched_user AS (
  SELECT 
    "u"."id",
    "ac"."id" as "authCodeId" -- Получаем ID кода аутентификации для последующего обновления
  FROM 
    "users" "u"
    INNER JOIN "auth_codes" "ac" ON "u"."id" = "ac"."userId"
  WHERE 
    "u"."email" = $1 -- Плейсхолдер для параметра email
    AND "ac"."code" = $2 -- Плейсхолдер для параметра кода аутентификации
    AND "ac"."updatedAt" >= NOW() - INTERVAL '10 minutes'
),
updated_users AS (
  UPDATE "users"
  SET "isVerified" = TRUE
  FROM matched_user
  WHERE "users"."id" = matched_user.id
  RETURNING "users"."id"
),
updated_auth_code AS (
  UPDATE "auth_codes"
  SET 
    "code" = floor(random() * 900000 + 100000), -- Генерируем новый код
    "updatedAt" = NOW() -- Обновляем время
  FROM matched_user
  WHERE "auth_codes"."id" = matched_user."authCodeId" -- Обновляем существующий код аутентификации
  RETURNING "auth_codes"."userId"
)
SELECT "id" FROM updated_users;