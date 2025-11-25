-- name: GetSnippetByID :one
SELECT snippet_id, user_id, title, content, created, expires
 	FROM snippets
 	WHERE expires > NOW()
 	AND snippet_id = $1;

-- name: CreateSnippet :one
INSERT INTO snippets (snippet_id, title, content, user_id, created, expires)
VALUES (gen_random_uuid(), $1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + CAST(@expires_in AS INTERVAL))
RETURNING snippet_id;

-- name: GetLatestSnippets :many
SELECT snippet_id, user_id, title, content, created, expires
 	FROM snippets
 	WHERE expires > NOW()
 	ORDER BY created DESC
 	LIMIT 10;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE snippet_id = $1;

-- name: CreateUser :one
INSERT INTO users (user_id, name, email, hashed_password)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING user_id;

-- name: GetUserByEmail :one
SELECT user_id, name, email, hashed_password
FROM users
WHERE email = $1;

-- name: GetLatestSnippetsForUser :many
SELECT snippet_id, title, content, created, expires
 	FROM snippets
 	WHERE expires > NOW() AND user_id = $1
 	ORDER BY created DESC
 	LIMIT 10;
