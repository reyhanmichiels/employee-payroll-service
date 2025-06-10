package attendance_period

import (
	"context"
	"testing"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	mock_auth "github.com/reyhanmichiels/go-pkg/v2/tests/mock/auth"
	mock_attendance_period "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_attendancePeriod_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockInterface(ctrl)
	mockAttendancePeriodDom := mock_attendance_period.NewMockInterface(ctrl)

	uc := Init(InitParam{
		Auth:             mockAuth,
		AttendancePeriod: mockAttendancePeriodDom,
	})

	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	defer func() { Now = time.Now }()

	mockLoginUser := auth.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockInputParam := dto.CreateAttendancePeriodParam{
		StartDate: null.DateFrom(mockTime.Add(24 * time.Hour)),
		EndDate:   null.DateFrom(mockTime.Add(48 * time.Hour)),
	}

	mockAttendancePeriod := entity.AttendancePeriod{
		ID:        1,
		StartDate: mockInputParam.StartDate,
		EndDate:   mockInputParam.EndDate,
		CreatedAt: null.TimeFrom(mockTime),
		CreatedBy: null.Int64From(mockLoginUser.ID),
	}

	tests := []struct {
		name       string
		inputParam dto.CreateAttendancePeriodParam
		mockFunc   func(inputParam dto.CreateAttendancePeriodParam)
		want       entity.AttendancePeriod
		wantErr    bool
	}{
		{
			name:       "Success",
			inputParam: mockInputParam,
			mockFunc: func(inputParam dto.CreateAttendancePeriodParam) {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Create(
					context.Background(),
					inputParam.ToAttendancePeriodInputParam(null.TimeFrom(mockTime), mockLoginUser.ID),
				).Return(mockAttendancePeriod, nil)
			},
			want:    mockAttendancePeriod,
			wantErr: false,
		},
		{
			name:       "Unique Constraint Error",
			inputParam: mockInputParam,
			mockFunc: func(inputParam dto.CreateAttendancePeriodParam) {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Create(
					context.Background(),
					inputParam.ToAttendancePeriodInputParam(null.TimeFrom(mockTime), mockLoginUser.ID),
				).Return(entity.AttendancePeriod{}, errors.NewWithCode(codes.CodeSQLUniqueConstraint, ""))
			},
			want:    entity.AttendancePeriod{},
			wantErr: true,
		},
		{
			name:       "Database Error",
			inputParam: mockInputParam,
			mockFunc: func(inputParam dto.CreateAttendancePeriodParam) {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Create(
					context.Background(),
					inputParam.ToAttendancePeriodInputParam(null.TimeFrom(mockTime), mockLoginUser.ID),
				).Return(entity.AttendancePeriod{}, assert.AnError)
			},
			want:    entity.AttendancePeriod{},
			wantErr: true,
		},
		{
			name:       "Validation Error",
			inputParam: dto.CreateAttendancePeriodParam{},
			mockFunc: func(inputParam dto.CreateAttendancePeriodParam) {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
			},
			want:    entity.AttendancePeriod{},
			wantErr: true,
		},
		{
			name:       "GetUserAuthInfo Error",
			inputParam: mockInputParam,
			mockFunc: func(inputParam dto.CreateAttendancePeriodParam) {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.User{}, assert.AnError)
			},
			want:    entity.AttendancePeriod{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.inputParam)

			got, err := uc.Create(context.Background(), tt.inputParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("attendancePeriod.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
