package attendance

const (
	insertAttendance = `
		INSERT INTO attendances (
			fk_attendance_period_id,
			fk_user_id,
			attendance_date,
			created_at,
			created_by
		) VALUES (
			:fk_attendance_period_id,
			:fk_user_id,
			:attendance_date,
			:created_at,
			:created_by
		) RETURNING *
	`

	readAttendance = `
		SELECT
			id,
			fk_attendance_period_id,
			fk_user_id,
			attendance_date,
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
			attendances
	`

	countAttendance = `
		SELECT
			COUNT(*)
		FROM
			attendances
	`

	updateAttendance = `
		UPDATE
			attendances
	`

	countUserAttendance = `
		SELECT
			fk_user_id AS userID,
			COUNT(*) AS attendanceCount
		FROM
		    attendances
		WHERE
		    fk_attendance_period_id = $1
		GROUP BY 
		    fk_user_id
	`
)
