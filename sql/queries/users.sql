-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1,$2)
RETURNING *;

-- name: VerifyUser :one
WITH matched_user AS (
  SELECT 
    u.id,
    ac.id as auth_code_id -- Получаем ID кода аутентификации для последующего обновления
  FROM 
    users u
    INNER JOIN auth_codes ac ON u.id = ac.user_id
  WHERE 
    u.email = $1 -- Плейсхолдер для параметра email
    AND ac.code = $2 -- Плейсхолдер для параметра кода аутентификации
    AND ac.created_at >= NOW() - INTERVAL '10 minutes'
),
updated_users AS (
  UPDATE users
  SET is_verified = TRUE
  FROM matched_user
  WHERE users.id = matched_user.id
  RETURNING users.id
),
updated_auth_code AS (
  UPDATE auth_codes
  SET 
    code = floor(random() * 900000 + 100000), -- Генерируем новый код
    updated_at = NOW() -- Обновляем время
  FROM matched_user
  WHERE auth_codes.id = matched_user.auth_code_id -- Обновляем существующий код аутентификации
  RETURNING auth_codes.user_id
)
SELECT id FROM updated_users;

-- name: GenerateAuthCode :one
WITH updated AS (
  UPDATE auth_codes
  SET code = floor(random() * 900000 + 100000)     
  FROM users  
  WHERE users.email = $1 
    AND auth_codes.user_id = users.id 
  RETURNING auth_codes.code
)
SELECT code FROM updated;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;




