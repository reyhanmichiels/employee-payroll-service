package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
)

// CreateAttendancePeriod godoc
// @Summary Create Attendance Period
// @Description Create a new attendance period
// @Tags Attendance Period
// @Security BearerAuth
// @Param data body dto.CreateAttendancePeriodParam true "Attendance Period Data"
// @Produce json
// @Success 201 {object} entity.HTTPResp{data=entity.AttendancePeriod{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 409 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/admin/attendance-periods [POST]
func (r *rest) CreateAttendancePeriod(ctx *gin.Context) {
	var param dto.CreateAttendancePeriodParam
	if err := r.Bind(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	data, err := r.uc.AttendancePeriod.Create(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeCreated, data, nil)
}

// GeneratePayroll godoc
// @Summary Generate Payroll
// @Description Generate payroll for a specific attendance period
// @Tags Attendance Period
// @Security BearerAuth
// @Param attendance_period_id path int true "Attendance Period ID"
// @Produce json
// @Success 202 {object} entity.HTTPResp{}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 409 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/admin/attendance-periods/{attendance_period_id}/payroll [POST]
func (r *rest) GeneratePayroll(ctx *gin.Context) {
	attendancePeriodIDStr := ctx.Param("attendance_period_id")
	if attendancePeriodIDStr == "" {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "attendance_period_id is empty"))
		return
	}

	attendancePeriodID, err := strconv.ParseInt(attendancePeriodIDStr, 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "attendance_period_id is not a valid number"))
		return
	}

	err = r.uc.AttendancePeriod.GeneratePayroll(ctx.Request.Context(), attendancePeriodID)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeAccepted, nil, nil)
}

// GeneratePayslip godoc
// @Summary Generate Payslip
// @Description Generate payslip for a specific attendance period
// @Tags Attendance Period
// @Security BearerAuth
// @Param attendance_period_id path int true "Attendance Period ID"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=dto.Payslip{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/attendance-periods/{attendance_period_id}/payslip [GET]
func (r *rest) GeneratePayslip(ctx *gin.Context) {
	attendancePeriodIDStr := ctx.Param("attendance_period_id")
	if attendancePeriodIDStr == "" {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "attendance_period_id is empty"))
		return
	}

	attendancePeriodID, err := strconv.ParseInt(attendancePeriodIDStr, 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "attendance_period_id is not a valid number"))
		return
	}

	data, err := r.uc.AttendancePeriod.GeneratePayslip(ctx.Request.Context(), attendancePeriodID)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, data, nil)
}
