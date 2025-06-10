package reimbursement

import (
	"context"
	"testing"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	mock_auth "github.com/reyhanmichiels/go-pkg/v2/tests/mock/auth"
	mock_reimbursement "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/reimbursement"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_reimbursement_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockInterface(ctrl)
	mockReimbursementDom := mock_reimbursement.NewMockInterface(ctrl)

	uc := Init(InitParam{
		Auth:          mockAuth,
		Reimbursement: mockReimbursementDom,
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

	mockInputParam := dto.CreateReimbursementParam{
		ReimbursementDate: null.DateFrom(mockTime.Add(-24 * time.Hour)),
		Amount:            100000,
		Description:       "Travel expenses",
	}

	mockReimbursementInputParam := mockInputParam.ToReimbursementInputParam(null.TimeFrom(mockTime), mockLoginUser.ID)
	mockReimbursementInputParam.MockApprovalData(null.TimeFrom(mockTime))

	mockReimbursement := entity.Reimbursement{
		ID: 1,
	}

	tests := []struct {
		name     string
		input    dto.CreateReimbursementParam
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Success",
			input: mockInputParam,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockReimbursementDom.EXPECT().Create(context.Background(), mockReimbursementInputParam).Return(mockReimbursement, nil)
			},
			wantErr: false,
		},
		{
			name:  "ReimbursementDom Create Error",
			input: mockInputParam,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(context.Background()).Return(mockLoginUser, nil)
				mockReimbursementDom.EXPECT().Create(context.Background(), mockReimbursementInputParam).Return(entity.Reimbursement{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:  "Validation Error",
			input: dto.CreateReimbursementParam{},
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
				t.Errorf("reimbursement.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
