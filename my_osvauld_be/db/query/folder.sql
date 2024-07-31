-- name: CreateFolder :one
INSERT INTO folders (
    name, description, created_by
) VALUES (
 $1, $2, $3
)
RETURNING id, created_at;



-- name: GetFoldersForUser :many
SELECT * FROM folders
WHERE created_by = $1;


