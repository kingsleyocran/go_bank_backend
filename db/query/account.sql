-- name: CreateAccount :one
INSERT INTO accounts (
  owner_name,
  balance,
  currency
) VALUES (
  $1, $2, $3
)RETURNING *;