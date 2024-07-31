-- name: AddFolderAccess :one
INSERT INTO folder_access (
    folder_id, user_id, access_type, group_id
) VALUES (
 $1, $2, $3, $4
)
RETURNING id, created_at;
