package attendance_period

import (
	"context"
	"fmt"
	"sync"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"golang.org/x/sync/errgroup"
)

func (a *attendancePeriod) PubSubGeneratePayroll(
	ctx context.Context,
	message entity.PubSubMessage,
) error {
	var (
		body dto.PubSubGeneratePayrollMessage
		err  error
	)

	defer func() {
		a.handleGeneratePayrollFailure(ctx, body.LoginUser.ID, body.AttendancePeriod.ID, err)
	}()

	if err = a.json.Unmarshal([]byte(message.Payload), &body); err != nil {
		return errors.NewWithCode(codes.CodeJSONUnmarshalError, "failed to unmarshal message body: %s", err.Error())
	}

	users, _, err := a.userDom.GetList(
		ctx,
		entity.UserParam{
			RoleID: entity.RoleIDUser,
			QueryOption: query.Option{
				IsActive:     true,
				DisableLimit: true,
			},
		},
	)
	if err != nil {
		return err
	}

	// Use goroutines to fetch data concurrently
	var (
		userAttendanceCount    map[int64]int64
		userIDToReimbursements map[int64][]entity.Reimbursement
		userIDToOvertimes      map[int64][]entity.Overtime
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		userAttendanceCount, err = a.attendanceDom.CountUserAttendance(gctx, body.AttendancePeriod.ID)
		return err
	})

	g.Go(func() error {
		var err error
		userIDToReimbursements, err = a.getUserIDToReimbursements(gctx, body.AttendancePeriod.StartDate, body.AttendancePeriod.EndDate)
		return err
	})

	g.Go(func() error {
		var err error
		userIDToOvertimes, err = a.getUserIDToOvertimes(gctx, body.AttendancePeriod.StartDate, body.AttendancePeriod.EndDate)
		return err
	})

	if err = g.Wait(); err != nil {
		return err
	}

	totalWorkingDays := body.AttendancePeriod.TotalWorkingDays()
	currentTime := null.TimeFrom(Now())
	userID := null.Int64From(body.LoginUser.ID)

	// Now implement goroutines for processing users
	err = a.transactor.Execute(ctx, "txPubSubGeneratePayroll", sql.TxOptions{}, func(ctx context.Context) error {
		const workerPoolSize = 5
		sem := make(chan struct{}, workerPoolSize)
		errChan := make(chan error, len(users))
		var wgUsers sync.WaitGroup

		for _, user := range users {
			wgUsers.Add(1)

			// Capture user variable for goroutine
			currentUser := user

			// Acquire semaphore slot
			sem <- struct{}{}

			go func() {
				defer wgUsers.Done()
				defer func() { <-sem }() // Release semaphore slot

				var payslipDetailInputParams []entity.PayslipDetailInputParam

				proratedSalary := currentUser.BaseSalary / float64(totalWorkingDays)

				basePayComponent, basePayDetail := a.calculateBasePayComponentAndDetail(
					userAttendanceCount[currentUser.ID],
					totalWorkingDays,
					proratedSalary,
				)

				overtimePayComponent, overtimeDetails := a.calculateOvertimePayComponentAndDetail(
					proratedSalary,
					userIDToOvertimes[currentUser.ID],
				)

				reimbursementPayComponent, reimbursementDetails := a.calculateReimbursementPayComponentAndDetail(
					proratedSalary,
					userIDToReimbursements[currentUser.ID],
				)

				payslip, err := a.payslipDom.Create(
					ctx,
					entity.PayslipInputParam{
						UserID:                 currentUser.ID,
						AttendancePeriodID:     body.AttendancePeriod.ID,
						BasePayComponent:       null.Float64From(basePayComponent),
						OvertimeComponent:      null.Float64From(overtimePayComponent),
						ReimbursementComponent: null.Float64From(reimbursementPayComponent),
						TotalTakeHomePay:       null.Float64From(basePayComponent + overtimePayComponent + reimbursementPayComponent),
						CreatedAt:              currentTime,
						CreatedBy:              userID,
					},
				)
				if err != nil {
					errChan <- err
					return
				}

				payslipDetailInputParams = append([]entity.PayslipDetailInputParam{basePayDetail}, overtimeDetails...)
				payslipDetailInputParams = append(payslipDetailInputParams, reimbursementDetails...)

				a.setPayslipIDToDetails(
					payslipDetailInputParams,
					payslip.ID,
					currentTime,
					userID,
				)

				err = a.payslipDetailDom.CreateMany(ctx, payslipDetailInputParams)
				if err != nil {
					errChan <- err
					return
				}
			}()
		}

		// Wait for all user processing to finish
		wgUsers.Wait()
		close(errChan)

		// Check if any errors occurred
		for err := range errChan {
			if err != nil {
				return err
			}
		}

		err := a.attendancePeriodDom.Update(
			ctx,
			entity.AttendancePeriodUpdateParam{
				PeriodStatus: entity.PeriodStatusProcessed,
				UpdatedAt:    null.TimeFrom(Now()),
				UpdatedBy:    null.Int64From(body.LoginUser.ID),
			},
			entity.AttendancePeriodParam{
				ID: body.AttendancePeriod.ID,
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *attendancePeriod) handleGeneratePayrollFailure(ctx context.Context, userID, attendancePeriodID int64, err error) {
	if err != nil {
		err := a.attendancePeriodDom.Update(
			ctx,
			entity.AttendancePeriodUpdateParam{
				PeriodStatus:        entity.PeriodStatusProcessError,
				PayrollProcessError: err.Error(),
				UpdatedAt:           null.TimeFrom(Now()),
				UpdatedBy:           null.Int64From(userID),
			},
			entity.AttendancePeriodParam{
				ID: attendancePeriodID,
			},
		)
		if err != nil {
			a.log.Error(ctx, fmt.Sprintf("failed to handle generate payroll failure: %s", err.Error()))
		}
	}
}

func (a *attendancePeriod) calculateBasePayComponentAndDetail(totalAttendance, totalWorkday int64, proratedSalary float64) (float64, entity.PayslipDetailInputParam) {
	totalPay := proratedSalary * float64(totalAttendance)
	return totalPay, entity.PayslipDetailInputParam{
		ItemType:    entity.PayslipItemTypeEarningBasePay,
		Description: fmt.Sprintf("Base Pay for %v Attendance on %v Workdays", totalAttendance, totalWorkday),
		Amount:      null.Float64From(totalPay),
	}
}

func (a *attendancePeriod) calculateOvertimePayComponentAndDetail(proratedSalary float64, overtimes []entity.Overtime) (float64, []entity.PayslipDetailInputParam) {
	inputParams := []entity.PayslipDetailInputParam{}
	totalHours := 0.0
	totalPay := 0.0

	for _, overtime := range overtimes {
		totalHours += overtime.OvertimeHour
		overtimePay := proratedSalary / 8 * overtime.OvertimeHour
		totalPay += overtimePay

		inputParams = append(inputParams, entity.PayslipDetailInputParam{
			ItemType:    entity.PayslipItemTypeEarningOvertime,
			Description: fmt.Sprintf("Overtime %v Hours on %v", overtime.OvertimeHour, overtime.OvertimeDate),
			Amount:      null.Float64From(overtimePay),
		})
	}

	return totalPay, inputParams
}

func (a *attendancePeriod) calculateReimbursementPayComponentAndDetail(proratedSalary float64, reimbursements []entity.Reimbursement) (float64, []entity.PayslipDetailInputParam) {
	inputParams := []entity.PayslipDetailInputParam{}
	totalPay := 0.0

	for _, reimbursement := range reimbursements {
		totalPay += reimbursement.Amount

		inputParams = append(inputParams, entity.PayslipDetailInputParam{
			ItemType:    entity.PayslipItemTypeReimbursement,
			Description: reimbursement.Description,
			Amount:      null.Float64From(reimbursement.Amount),
		})
	}

	return totalPay, inputParams
}

func (a *attendancePeriod) getUserIDToReimbursements(
	ctx context.Context,
	startDate null.Date,
	endDate null.Date,
) (map[int64][]entity.Reimbursement, error) {
	userIDToReimbursements := make(map[int64][]entity.Reimbursement)

	reimbursements, _, err := a.reimbursementDom.GetList(
		ctx,
		entity.ReimbursementParam{
			ApprovedDateGTE: startDate,
			ApprovedDateLTE: endDate,
			QueryOption: query.Option{
				IsActive:     true,
				DisableLimit: true,
			},
		},
	)
	if err != nil {
		return userIDToReimbursements, err
	}

	for _, reimbursement := range reimbursements {
		userIDToReimbursements[reimbursement.UserID] = append(userIDToReimbursements[reimbursement.UserID], reimbursement)
	}

	return userIDToReimbursements, nil
}

func (a *attendancePeriod) getUserIDToOvertimes(
	ctx context.Context,
	startDate null.Date,
	endDate null.Date,
) (map[int64][]entity.Overtime, error) {
	userIDToOvertimes := make(map[int64][]entity.Overtime)

	overtimes, _, err := a.overtimeDom.GetList(
		ctx,
		entity.OvertimeParam{
			ApprovedDateGTE: startDate,
			ApprovedDateLTE: endDate,
			QueryOption: query.Option{
				IsActive:     true,
				DisableLimit: true,
			},
		},
	)
	if err != nil {
		return userIDToOvertimes, err
	}

	for _, overtime := range overtimes {
		userIDToOvertimes[overtime.UserID] = append(userIDToOvertimes[overtime.UserID], overtime)
	}

	return userIDToOvertimes, nil
}

func (a *attendancePeriod) setPayslipIDToDetails(
	details []entity.PayslipDetailInputParam,
	payslipID int64,
	currentTime null.Time,
	userID null.Int64,
) {
	for i := range details {
		details[i].PayslipID = payslipID
		details[i].CreatedAt = currentTime
		details[i].CreatedBy = userID
	}
}
