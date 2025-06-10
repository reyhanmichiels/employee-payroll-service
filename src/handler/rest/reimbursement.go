package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
)

// SubmitReimbursement godoc
// @Summary Submit Reimbursement
// @Description Submit reimbursement for a user
// @Tags Reimbursement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param param body dto.CreateReimbursementParam true "Reimbursement Request Parameters"
// @Success 201 {object} entity.HTTPResp{data=entity.Reimbursement{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/reimbursements [POST]
func (r *rest) SubmitReimbursement(ctx *gin.Context) {
	var param dto.CreateReimbursementParam
	if err := r.Bind(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	data, err := r.uc.Reimbursement.Create(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeCreated, data, nil)
}
