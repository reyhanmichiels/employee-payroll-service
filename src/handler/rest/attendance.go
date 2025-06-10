package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
)

// SubmitAttendance godoc
// @Summary Submit Attendance
// @Description Submit attendance for a user
// @Tags Attendance
// @Security BearerAuth
// @Produce json
// @Success 201 {object} entity.HTTPResp{}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 409 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/attendances [POST]
func (r *rest) SubmitAttendance(ctx *gin.Context) {
	err := r.uc.Attendance.Create(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeCreated, nil, nil)
}
