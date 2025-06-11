package payslip

const (
	insertPayslip = `
		INSERT INTO payslips (
			fk_user_id,
			fk_attendance_period_id,
			base_pay_component,
			overtime_component,
			reimbursement_component,
			total_take_home_pay,
			created_at,
			created_by
		) VALUES (
			:fk_user_id,
			:fk_attendance_period_id,
			:base_pay_component,
			:overtime_component,
			:reimbursement_component,
			:total_take_home_pay,
			:created_at,
			:created_by
		) RETURNING *
	`

	insertManyPayslip = `
		INSERT INTO payslips (
			fk_user_id,
			fk_attendance_period_id,
			base_pay_component,
			overtime_component,
			reimbursement_component,
			total_take_home_pay,
			created_at,
			created_by
		) VALUES (
			:fk_user_id,
			:fk_attendance_period_id,
			:base_pay_component,
			:overtime_component,
			:reimbursement_component,
			:total_take_home_pay,
			:created_at,
			:created_by
		)
	`

	readPayslip = `
		SELECT
			id,
			fk_user_id,
			fk_attendance_period_id,
			base_pay_component,
			overtime_component,
			reimbursement_component,
			total_take_home_pay,
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
			payslips
	`

	countPayslip = `
		SELECT
			COUNT(*)
		FROM
			payslips
	`

	updatePayslip = `
		UPDATE
			payslips
	`
)
