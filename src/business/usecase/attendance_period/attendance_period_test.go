package attendance_period

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
	mock_parser "github.com/reyhanmichiels/go-pkg/v2/tests/mock/parser"
	mock_attendance "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/attendance"
	mock_attendance_period "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/attendance_period"
	mock_overtime "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/overtime"
	mock_payslip "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/payslip"
	mock_payslip_detail "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/payslip_detail"
	mock_reimbursement "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/reimbursement"
	mock_transactor "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/transactor"
	mock_user "github.com/reyhanmichies/employee-payroll-service/src/business/domain/mock/user"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	mock_publisher "github.com/reyhanmichies/employee-payroll-service/src/handler/pubsub/mock/publisher"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockField struct {
	auth                *mock_auth.MockInterface
	attendancePeriodDom *mock_attendance_period.MockInterface
	publisher           *mock_publisher.MockInterface
	transactor          *mock_transactor.MockInterface
	attendanceDom       *mock_attendance.MockInterface
	userDom             *mock_user.MockInterface
	reimbursementDom    *mock_reimbursement.MockInterface
	overtimeDom         *mock_overtime.MockInterface
	payslipDom          *mock_payslip.MockInterface
	payslipDetailDom    *mock_payslip_detail.MockInterface
	json                *mock_parser.MockJSONInterface
}

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

func Test_attendancePeriod_GeneratePayroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockInterface(ctrl)
	mockAttendancePeriodDom := mock_attendance_period.NewMockInterface(ctrl)
	mockPublisher := mock_publisher.NewMockInterface(ctrl)
	mockTransactor := mock_transactor.NewMockInterface(ctrl)

	uc := Init(InitParam{
		Auth:             mockAuth,
		AttendancePeriod: mockAttendancePeriodDom,
		Publisher:        mockPublisher,
		Transactor:       mockTransactor,
	})

	mock := mockField{
		auth:                mockAuth,
		attendancePeriodDom: mockAttendancePeriodDom,
		publisher:           mockPublisher,
		transactor:          mockTransactor,
	}

	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	restoreAll := func() {
		Now = time.Now
	}
	defer restoreAll()

	type args struct {
		ctx                context.Context
		attendancePeriodID int64
	}

	mockLoginUser := auth.User{
		ID:     100,
		Name:   "Test User",
		Email:  "test@example.com",
		RoleID: 1,
	}

	mockAttendancePeriodParam := entity.AttendancePeriodParam{
		ID: 1,
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	mockAttendancePeriodUpdateParam := entity.AttendancePeriodUpdateParam{
		PeriodStatus: entity.PeriodStatusProcessed,
		UpdatedAt:    null.TimeFrom(mockTime),
		UpdatedBy:    null.Int64From(mockLoginUser.ID),
	}

	mockAttendancePeriod := entity.AttendancePeriod{
		ID:           1,
		StartDate:    null.DateFrom(time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:      null.DateFrom(time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC)),
		PeriodStatus: entity.PeriodStatusClosed,
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockField, _args args)
		wantErr  bool
	}{
		{
			name: "Success",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(mockAttendancePeriod, nil)
				mock.transactor.
					EXPECT().
					Execute(
						_args.ctx,
						"txGeneratePayroll",
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(
					func(_ context.Context, _ string, _ interface{}, callback func(context.Context) error) error {
						return callback(_args.ctx)
					},
				)
				mock.attendancePeriodDom.EXPECT().Update(_args.ctx, mockAttendancePeriodUpdateParam, mockAttendancePeriodParam).Return(nil)
				mock.publisher.EXPECT().Publish(
					_args.ctx,
					entity.ExchangePayrollEvent,
					entity.RoutingKeyPayrollCalculate,
					dto.PubSubGeneratePayrollMessage{
						AttendancePeriod: mockAttendancePeriod,
						LoginUser:        mockLoginUser,
					},
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Publish",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(mockAttendancePeriod, nil)
				mock.transactor.
					EXPECT().
					Execute(
						_args.ctx,
						"txGeneratePayroll",
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(
					func(_ context.Context, _ string, _ interface{}, callback func(context.Context) error) error {
						return callback(_args.ctx)
					},
				)
				mock.attendancePeriodDom.EXPECT().Update(_args.ctx, mockAttendancePeriodUpdateParam, mockAttendancePeriodParam).Return(nil)
				mock.publisher.EXPECT().Publish(
					_args.ctx,
					entity.ExchangePayrollEvent,
					entity.RoutingKeyPayrollCalculate,
					dto.PubSubGeneratePayrollMessage{
						AttendancePeriod: mockAttendancePeriod,
						LoginUser:        mockLoginUser,
					},
				).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed Update",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(mockAttendancePeriod, nil)
				mock.transactor.
					EXPECT().
					Execute(
						_args.ctx,
						"txGeneratePayroll",
						gomock.Any(),
						gomock.Any(),
					).DoAndReturn(
					func(_ context.Context, _ string, _ interface{}, callback func(context.Context) error) error {
						return callback(_args.ctx)
					},
				)
				mock.attendancePeriodDom.EXPECT().Update(_args.ctx, mockAttendancePeriodUpdateParam, mockAttendancePeriodParam).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed Transactor",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(mockAttendancePeriod, nil)
				mock.transactor.
					EXPECT().
					Execute(
						_args.ctx,
						"txGeneratePayroll",
						gomock.Any(),
						gomock.Any(),
					).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed Status Processed",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(
					entity.AttendancePeriod{ID: 1, PeriodStatus: entity.PeriodStatusProcessed}, nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed Status Processing",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(
					entity.AttendancePeriod{ID: 1, PeriodStatus: entity.PeriodStatusProcessing}, nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed Status Open",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(
					entity.AttendancePeriod{ID: 1, PeriodStatus: entity.PeriodStatusOpen}, nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed Status Upcoming",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(
					entity.AttendancePeriod{ID: 1, PeriodStatus: entity.PeriodStatusUpcoming}, nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed Get AttendancePeriod Not Found",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(
					entity.AttendancePeriod{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "not found"),
				)
			},
			wantErr: true,
		},
		{
			name: "Failed Get AttendancePeriod Generic Error",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(mockLoginUser, nil)
				mock.attendancePeriodDom.EXPECT().Get(_args.ctx, mockAttendancePeriodParam).Return(
					entity.AttendancePeriod{}, assert.AnError,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed Get UserAuthInfo",
			args: args{
				ctx:                context.Background(),
				attendancePeriodID: 1,
			},
			mockFunc: func(mock mockField, _args args) {
				mock.auth.EXPECT().GetUserAuthInfo(_args.ctx).Return(auth.User{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock, tt.args)

			err := uc.GeneratePayroll(tt.args.ctx, tt.args.attendancePeriodID)
			if (err != nil) != tt.wantErr {
				t.Errorf("attendancePeriod.GeneratePayroll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_attendancePeriod_GeneratePayslip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockInterface(ctrl)
	mockAttendancePeriodDom := mock_attendance_period.NewMockInterface(ctrl)
	mockPayslipDom := mock_payslip.NewMockInterface(ctrl)
	mockPayslipDetailDom := mock_payslip_detail.NewMockInterface(ctrl)

	uc := Init(InitParam{
		Auth:             mockAuth,
		AttendancePeriod: mockAttendancePeriodDom,
		Payslip:          mockPayslipDom,
		PayslipDetail:    mockPayslipDetailDom,
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

	mockAttendancePeriodID := int64(1)
	mockAttendancePeriod := entity.AttendancePeriod{
		ID:        mockAttendancePeriodID,
		StartDate: null.DateFrom(time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   null.DateFrom(time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC)),
	}

	mockPayslip := entity.Payslip{
		ID:                     1,
		UserID:                 mockLoginUser.ID,
		AttendancePeriodID:     mockAttendancePeriodID,
		BasePayComponent:       1000.0,
		OvertimeComponent:      200.0,
		ReimbursementComponent: 100.0,
		TotalTakeHomePay:       1300.0,
	}

	mockPayslipDetails := []entity.PayslipDetail{
		{
			ID:          1,
			PayslipID:   mockPayslip.ID,
			ItemType:    entity.PayslipItemTypeEarningBasePay,
			Description: "Base Pay for 20 Attendance on 22 Workdays",
			Amount:      1000.0,
		},
		{
			ID:          2,
			PayslipID:   mockPayslip.ID,
			ItemType:    entity.PayslipItemTypeEarningOvertime,
			Description: "Overtime 4 Hours on 2023-05-10",
			Amount:      200.0,
		},
		{
			ID:          3,
			PayslipID:   mockPayslip.ID,
			ItemType:    entity.PayslipItemTypeReimbursement,
			Description: "Transport Reimbursement",
			Amount:      100.0,
		},
	}

	expectedPayslip := dto.Payslip{
		StartDate:          mockAttendancePeriod.StartDate,
		EndDate:            mockAttendancePeriod.EndDate,
		BasePayComponent:   mockPayslip.BasePayComponent,
		OvertimeComponent:  mockPayslip.OvertimeComponent,
		ReimburseComponent: mockPayslip.ReimbursementComponent,
		TotalTakeHomePay:   mockPayslip.TotalTakeHomePay,
		Details: []dto.PayslipDetail{
			{
				Type:        entity.PayslipItemTypeEarningBasePay,
				Description: "Base Pay for 20 Attendance on 22 Workdays",
				Amount:      1000.0,
			},
			{
				Type:        entity.PayslipItemTypeEarningOvertime,
				Description: "Overtime 4 Hours on 2023-05-10",
				Amount:      200.0,
			},
			{
				Type:        entity.PayslipItemTypeReimbursement,
				Description: "Transport Reimbursement",
				Amount:      100.0,
			},
		},
	}

	tests := []struct {
		name               string
		attendancePeriodID int64
		mockFunc           func()
		want               dto.Payslip
		wantErr            bool
	}{
		{
			name:               "Success",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(
					gomock.Any(),
					entity.AttendancePeriodParam{
						ID: mockAttendancePeriodID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(mockAttendancePeriod, nil)
				mockPayslipDom.EXPECT().Get(
					gomock.Any(),
					entity.PayslipParam{
						AttendancePeriodID: mockAttendancePeriodID,
						UserID:             mockLoginUser.ID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(mockPayslip, nil)
				mockPayslipDetailDom.EXPECT().GetList(
					gomock.Any(),
					entity.PayslipDetailParam{
						PayslipID: mockPayslip.ID,
						QueryOption: query.Option{
							IsActive:     true,
							DisableLimit: true,
						},
					},
				).Return(mockPayslipDetails, nil, nil)
			},
			want:    expectedPayslip,
			wantErr: false,
		},
		{
			name:               "Failed PayslipDetail Generic Error",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(
					gomock.Any(),
					entity.AttendancePeriodParam{
						ID: mockAttendancePeriodID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(mockAttendancePeriod, nil)
				mockPayslipDom.EXPECT().Get(
					gomock.Any(),
					entity.PayslipParam{
						AttendancePeriodID: mockAttendancePeriodID,
						UserID:             mockLoginUser.ID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(mockPayslip, nil)
				mockPayslipDetailDom.EXPECT().GetList(
					gomock.Any(),
					entity.PayslipDetailParam{
						PayslipID: mockPayslip.ID,
						QueryOption: query.Option{
							IsActive:     true,
							DisableLimit: true,
						},
					},
				).Return(nil, nil, assert.AnError)
			},
			want:    dto.Payslip{},
			wantErr: true,
		},
		{
			name:               "Failed Payslip Not Found",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(
					gomock.Any(),
					entity.AttendancePeriodParam{
						ID: mockAttendancePeriodID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(mockAttendancePeriod, nil)
				mockPayslipDom.EXPECT().Get(
					gomock.Any(),
					entity.PayslipParam{
						AttendancePeriodID: mockAttendancePeriodID,
						UserID:             mockLoginUser.ID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(entity.Payslip{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "not found"))
			},
			want:    dto.Payslip{},
			wantErr: true,
		},
		{
			name:               "Failed Payslip Generic Error",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(
					gomock.Any(),
					entity.AttendancePeriodParam{
						ID: mockAttendancePeriodID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(mockAttendancePeriod, nil)
				mockPayslipDom.EXPECT().Get(
					gomock.Any(),
					entity.PayslipParam{
						AttendancePeriodID: mockAttendancePeriodID,
						UserID:             mockLoginUser.ID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(entity.Payslip{}, assert.AnError)
			},
			want:    dto.Payslip{},
			wantErr: true,
		},
		{
			name:               "Failed AttendancePeriod Not Found",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(
					gomock.Any(),
					entity.AttendancePeriodParam{
						ID: mockAttendancePeriodID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(entity.AttendancePeriod{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "not found"))
			},
			want:    dto.Payslip{},
			wantErr: true,
		},
		{
			name:               "Failed AttendancePeriod Generic Error",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockLoginUser, nil)
				mockAttendancePeriodDom.EXPECT().Get(
					gomock.Any(),
					entity.AttendancePeriodParam{
						ID: mockAttendancePeriodID,
						QueryOption: query.Option{
							IsActive: true,
						},
					},
				).Return(entity.AttendancePeriod{}, assert.AnError)
			},
			want:    dto.Payslip{},
			wantErr: true,
		},
		{
			name:               "Failed GetUserAuthInfo",
			attendancePeriodID: mockAttendancePeriodID,
			mockFunc: func() {
				mockAuth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(auth.User{}, assert.AnError)
			},
			want:    dto.Payslip{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			got, err := uc.GeneratePayslip(context.Background(), tt.attendancePeriodID)
			if (err != nil) != tt.wantErr {
				t.Errorf("attendancePeriod.GeneratePayslip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
