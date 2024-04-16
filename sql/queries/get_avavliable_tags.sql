-- name: GetAvailableTags :many
SELECT ut.*
FROM user_tags ut
WHERE NOT EXISTS (
    SELECT 1
    FROM user_user_tags uut
    WHERE uut."userTagId" = ut.id
    AND uut."userId" = $1
);
