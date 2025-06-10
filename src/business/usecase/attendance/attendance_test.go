package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	mock_auth "github.com/reyhanmichiels/go-pkg/v2/tests/mock/auth"
	attendance_dom "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/attendance"
	attendance_period_dom "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_attendance_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockInterface(ctrl)
	mockAttendanceDom := attendance_dom.NewMockInterface(ctrl)
	mockAttendancePeriodDom := attendance_period_dom.NewMockInterface(ctrl)

	uc := Init(InitParam{
		AttendancePeriod: mockAttendancePeriodDom,
		Attendance:       mockAttendanceDom,
		Auth:             mockAuth,
	})

	mockTime := time.Date(2023, 10, 6, 10, 0, 0, 0, time.UTC) // A Friday
	Now = func() time.Time {
		return mockTime
	}
	defer func() { Now = time.Now }()

	mockLoginUser := auth.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockAttendancePeriod := entity.AttendancePeriod{
		ID: 1,
	}

	tests := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(gomock.Any(), entity.AttendancePeriodParam{
					PeriodStatus: entity.PeriodStatusOpen,
					QueryOption: query.Option{
						IsActive: true,
					},
				}).Return(mockAttendancePeriod, nil)
				mockAttendanceDom.EXPECT().Create(gomock.Any(), entity.AttendanceInputParam{
					AttendancePeriodID: mockAttendancePeriod.ID,
					UserID:             mockLoginUser.ID,
					AttendanceDate:     null.DateFrom(mockTime),
					CreatedAt:          null.TimeFrom(mockTime),
					CreatedBy:          null.Int64From(mockLoginUser.ID),
				}).Return(entity.Attendance{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Attendance Already Submitted",
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(gomock.Any(), entity.AttendancePeriodParam{
					PeriodStatus: entity.PeriodStatusOpen,
					QueryOption: query.Option{
						IsActive: true,
					},
				}).Return(mockAttendancePeriod, nil)
				mockAttendanceDom.EXPECT().Create(gomock.Any(), entity.AttendanceInputParam{
					AttendancePeriodID: mockAttendancePeriod.ID,
					UserID:             mockLoginUser.ID,
					AttendanceDate:     null.DateFrom(mockTime),
					CreatedAt:          null.TimeFrom(mockTime),
					CreatedBy:          null.Int64From(mockLoginUser.ID),
				}).Return(entity.Attendance{}, errors.NewWithCode(codes.CodeSQLUniqueConstraint, ""))
			},
			wantErr: true,
		},
		{
			name: "Database Error When Submit Attendance",
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(gomock.Any(), entity.AttendancePeriodParam{
					PeriodStatus: entity.PeriodStatusOpen,
					QueryOption: query.Option{
						IsActive: true,
					},
				}).Return(mockAttendancePeriod, nil)
				mockAttendanceDom.EXPECT().Create(gomock.Any(), entity.AttendanceInputParam{
					AttendancePeriodID: mockAttendancePeriod.ID,
					UserID:             mockLoginUser.ID,
					AttendanceDate:     null.DateFrom(mockTime),
					CreatedAt:          null.TimeFrom(mockTime),
					CreatedBy:          null.Int64From(mockLoginUser.ID),
				}).Return(entity.Attendance{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "No Open Attendance Period",
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(gomock.Any(), entity.AttendancePeriodParam{
					PeriodStatus: entity.PeriodStatusOpen,
					QueryOption: query.Option{
						IsActive: true,
					},
				}).Return(entity.AttendancePeriod{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, ""))
			},
			wantErr: true,
		},
		{
			name: "Database Error When Get Attendance Period",
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(gomock.Any(), entity.AttendancePeriodParam{
					PeriodStatus: entity.PeriodStatusOpen,
					QueryOption: query.Option{
						IsActive: true,
					},
				}).Return(entity.AttendancePeriod{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Weekend Error",
			mockFunc: func() {
				Now = func() time.Time {
					return time.Date(2023, 10, 7, 10, 0, 0, 0, time.UTC) // A Saturday
				}
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(auth.User{}, nil)
			},
			wantErr: true,
		},
		{
			name: "GetUserAuthInfo Error",
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(auth.User{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := uc.Create(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("attendance.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
