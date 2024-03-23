-- name: CreateAuthCode :one
INSERT INTO auth_codes (user_id)
VALUES ($1)
RETURNING code;


-- name: UpdateAuthCode :one
UPDATE auth_codes
SET code = $1  -- Replace $1 with the new code value
WHERE user_id = $2  -- Replace $2 with the user ID
RETURNING code;  -- Optional: retrieve the updated code


-- name: CheckAuthCode :one
WITH matched_user AS (
    SELECT auth_codes.user_id
    FROM auth_codes
    JOIN users ON auth_codes.user_id = users.id
    WHERE users.email = $1
      AND users.is_verified=true
      AND auth_codes.code = $2
      AND auth_codes.updated_at > (CURRENT_TIMESTAMP - INTERVAL '10 MINUTES')
    LIMIT 1
),
updated AS (
    UPDATE auth_codes
    SET code = floor(random() * 900000 + 100000) -- Generates a new code        
    WHERE user_id IN (SELECT user_id FROM matched_user)
    RETURNING user_id
)
SELECT user_id FROM updated;