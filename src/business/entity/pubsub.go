package entity

// exchange name
const (
	ExchangePayrollEvent = "payroll.event"
)

// queue name
const (
	QueuePayrollCalculation = "employee-payroll-service.payroll-calculation"
)

// routing key
const (
	RoutingKeyPayrollCalculate = "payroll.calculate"
)

type PubSubMessage struct {
	RequestID string `json:"requestId"`
	Event     string `json:"event"`
	Payload   string `json:"payload"`
}
