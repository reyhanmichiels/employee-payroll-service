package payslip_detail

const (
	insertPayslipDetail = `
		INSERT INTO payslip_details (
			fk_payslip_id,
			item_type,
			description,
			amount,
			created_at,
			created_by
		) VALUES (
			:fk_payslip_id,
			:item_type,
			:description,
			:amount,
			:created_at,
			:created_by
		) RETURNING *
	`

	insertManyPayslipDetail = `
		INSERT INTO payslip_details (
			fk_payslip_id,
			item_type,
			description,
			amount,
			created_at,
			created_by
		) VALUES (
			:fk_payslip_id,
			:item_type,
			:description,
			:amount,
			:created_at,
			:created_by
		)
	`

	readPayslipDetail = `
		SELECT
			id,
			fk_payslip_id,
			item_type,
			description,
			amount,
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
			payslip_details
	`

	countPayslipDetail = `
		SELECT
			COUNT(*)
		FROM
			payslip_details
	`

	updatePayslipDetail = `
		UPDATE
			payslip_details
	`
)
