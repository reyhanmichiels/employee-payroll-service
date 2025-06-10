package overtime

import (
	"context"
	"testing"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	mock_auth "github.com/reyhanmichiels/go-pkg/v2/tests/mock/auth"
	mock_overtime "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/overtime"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_overtime_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockInterface(ctrl)
	mockOvertimeDom := mock_overtime.NewMockInterface(ctrl)

	uc := Init(InitParam{
		Auth:        mockAuth,
		OvertimeDom: mockOvertimeDom,
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

	mockInputParam := dto.CreateOvertimeParam{
		OvertimeDate: null.DateFrom(mockTime.Add(-24 * time.Hour)),
		OvertimeHour: null.Float64From(1.5),
	}

	mockOvertimeInputParam := mockInputParam.ToOvertimeInputParam(null.TimeFrom(mockTime), mockLoginUser.ID)
	mockOvertimeInputParam.MockApprovalData(null.TimeFrom(mockTime))

	mockOvertime := entity.Overtime{
		ID: 1,
	}

	tests := []struct {
		name     string
		input    dto.CreateOvertimeParam
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Success",
			input: mockInputParam,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockOvertimeDom.EXPECT().Create(context.Background(), mockOvertimeInputParam).Return(mockOvertime, nil)
			},
			wantErr: false,
		},
		{
			name:  "Duplicate Overtime Submission",
			input: mockInputParam,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockOvertimeDom.EXPECT().Create(context.Background(), mockOvertimeInputParam).Return(entity.Overtime{}, errors.NewWithCode(codes.CodeSQLUniqueConstraint, ""))
			},
			wantErr: true,
		},
		{
			name:  "OvertimeDom Create Error",
			input: mockInputParam,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockOvertimeDom.EXPECT().Create(context.Background(), mockOvertimeInputParam).Return(entity.Overtime{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:  "Validation Error",
			input: dto.CreateOvertimeParam{},
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
			},
			wantErr: true,
		},
		{
			name:  "GetUserAuthInfo Error",
			input: mockInputParam,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.User{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			_, err := uc.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("overtime.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
