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