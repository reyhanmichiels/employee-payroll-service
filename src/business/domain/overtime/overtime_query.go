package overtime

const (
	insertOvertime = `
		INSERT INTO overtimes (
			fk_user_id,
			overtime_date,
			overtime_hour,
			approved_date,
			approved_by,
			created_at,
			created_by
		) VALUES (
			:fk_user_id,
			:overtime_date,
			:overtime_hour,
			:approved_date,
			:approved_by,
			:created_at,
			:created_by
		) RETURNING *
	`

	readOvertime = `
		SELECT
			id,
			fk_user_id,
			overtime_date,
			overtime_hour,
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
			overtimes
	`

	countOvertime = `
		SELECT
			COUNT(*)
		FROM
			overtimes
	`

	updateOvertime = `
		UPDATE
			overtimes
	`
)
