-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;


-- name: UpdateAccountsBalance :many
WITH deduct AS (
  UPDATE accounts
  SET balance = accounts.balance - sqlc.arg(amount)
  WHERE accounts.id = sqlc.arg(fromAccountID) AND balance >= sqlc.arg(amount)
  RETURNING *
),
add AS (
  UPDATE accounts
  SET balance = accounts.balance + sqlc.arg(amount)
  WHERE accounts.id = sqlc.arg(ToAccountID)
  RETURNING *
)
SELECT * FROM deduct
UNION ALL
SELECT * FROM add;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY accounts.id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
  set balance = $2
WHERE id = $1 RETURNING *;

-- name: DeleteAccound :exec
DELETE FROM accounts
WHERE id = $1;