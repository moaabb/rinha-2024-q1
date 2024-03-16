package transactiondb

const GetAccountStatementByPartyId = `
	SELECT t.transaction_id, t.value, t.type, t.description, t.created_at
	FROM rinha.party p
	INNER JOIN rinha.transactions t
	ON p.party_id = t.party_id
	WHERE p.party_id = $1
	ORDER BY created_at DESC
	LIMIT 10
`

const CreateTransaction = `
	INSERT INTO rinha.transactions
		(value, type, description, created_at, party_id)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING value, type, description, created_at, party_id
`

const AddToAccountBalance = `
	UPDATE rinha.party SET balance = balance + $1  WHERE party_id = $2 RETURNING party_id, "limit", balance
`

const SubtractFromAccountBalance = `
	UPDATE rinha.party SET balance = balance - $1  WHERE party_id = $2 RETURNING party_id, "limit", balance
`

const UpdateAccountBalance = `
	UPDATE rinha.party SET balance =  $1  WHERE party_id = $2 RETURNING party_id, "limit", balance
`

const GetAccountBalance = `
	SELECT party_id, "limit", balance FROM rinha.party WHERE party_id = $1
`
