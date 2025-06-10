package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
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
// @Router /v1/attendance-periods [POST]
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
