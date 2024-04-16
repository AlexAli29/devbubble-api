-- name: GetUserById :one
SELECT 
    u.id, 
    u."createdAt" , 
    u.email, 
    u.description, 
    u.name, 
    u."isVerified",
    json_agg(json_build_object('id', ut.id, 'text', ut.text, 'icon', ut.icon)) FILTER (WHERE ut.id IS NOT NULL) AS tags
FROM users u
LEFT JOIN user_user_tags uut ON u.id = uut."userId"
LEFT JOIN user_tags ut ON uut."userTagId" = ut.id
WHERE u.id = $1
GROUP BY u.id;


-- name: GetUserIsVerifiedEmailNameById :one
SELECT "isVerified", "email", "name" FROM "users" WHERE "id" = $1;