-- name: GetById :one
SELECT snippet_id, title, content, created, expires
 	FROM snippets
 	WHERE expires > NOW()
 	AND snippet_id = $1;

-- name: Create :one
INSERT INTO snippets (snippet_id, title, content, created, expires)
VALUES (gen_random_uuid(), $1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + CAST(@expires_in AS INTERVAL))
RETURNING snippet_id;

-- name: GetLatest :many
SELECT snippet_id, title, content, created, expires
 	FROM snippets
 	WHERE expires > NOW()
 	ORDER BY created DESC
 	LIMIT 10;