package transactiondb

const GetAccountStatementByPartyId = `
	SELECT transaction_id, value, type, description, created_at
	FROM rinha.party p
	INNER JOIN
	(
		SELECT
			party_id, value, description, created_at, transaction_id, type
		FROM rinha.transactions t
		ORDER BY created_at DESC
		LIMIT 10
	) t
	ON p.party_id = t.party_id
	WHERE p.party_id = $1
`

const CreateTransaction = `
	INSERT INTO rinha.transactions
		(value, type, description, created_at, party_id)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING value, type, description, created_at, party_id
`

const ChangeAccountBalance = `
	UPDATE rinha.party SET balance = $1 WHERE party_id = $2 RETURNING party_id, "limit", balance
`

const GetAccountBalance = `
	SELECT party_id, "limit", balance FROM rinha.party WHERE party_id = $1
`
