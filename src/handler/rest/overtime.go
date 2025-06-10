package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
)

// SubmitOvertime godoc
// @Summary Submit Overtime
// @Description Submit overtime for a user
// @Tags Overtime
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param param body dto.CreateOvertimeParam true "Overtime Request Parameters"
// @Success 201 {object} entity.HTTPResp{data=entity.Overtime{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 409 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/overtimes [POST]
func (r *rest) SubmitOvertime(ctx *gin.Context) {
	var param dto.CreateOvertimeParam
	if err := r.Bind(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	data, err := r.uc.Overtime.Create(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeCreated, data, nil)
}
