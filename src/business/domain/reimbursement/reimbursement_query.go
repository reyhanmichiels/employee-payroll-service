package reimbursement

const (
	insertReimbursement = `
		INSERT INTO reimbursements (
			fk_user_id,
			description,
			amount,
			approved_date,
			approved_by,
			created_at,
			created_by
		) VALUES (
			:fk_user_id,
			:description,
			:amount,
			:approved_date,
			:approved_by,
			:created_at,
			:created_by
		)
	`

	readReimbursement = `
		SELECT
			id,
			fk_user_id,
			description,
			amount,
			approved_date,
			approved_by,
			status,
			flag,
			meta,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			reimbursements
	`

	countReimbursement = `
		SELECT
			COUNT(*)
		FROM
			reimbursements
	`

	updateReimbursement = `
		UPDATE
			reimbursements
	`
)
