package attendance_period

const (
	insertAttendancePeriod = `
		INSERT INTO attendance_periods (
			start_date,
			end_date,
			period_status,
			created_at,
			created_by
		) VALUES (
			:start_date,
			:end_date,
			:period_status,
			:created_at,
			:created_by
		)
	`

	readAttendancePeriod = `
		SELECT
			id,
			start_date,
			end_date,
			period_status,
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
			attendance_periods
	`

	countAttendancePeriod = `
		SELECT
			COUNT(*)
		FROM
			attendance_periods
	`

	updateAttendancePeriod = `
		UPDATE
			attendance_periods
	`
)
