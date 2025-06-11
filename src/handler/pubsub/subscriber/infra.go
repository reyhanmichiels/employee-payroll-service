package subscriber

import (
	"context"
	"fmt"

	"github.com/reyhanmichiels/go-pkg/v2/pubsub/rabbitmq"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/reyhanmichies/employee-payroll-service/src/handler/pubsub"
)

type Exchange struct {
	Name string
	Type rabbitmq.ExchangeType
}

type Queue struct {
	Name string
}

type QueueBind struct {
	ExchangeName string
	QueueName    string
	RoutingKey   string
}

var MQExchanges = []Exchange{
	{
		Name: entity.ExchangePayrollEvent,
		Type: rabbitmq.ExchangeTypeTopic,
	},
}

var MQQueue = []Queue{
	{
		Name: entity.QueuePayrollCalculation,
	},
}

var MQQueueBind = []QueueBind{
	{
		ExchangeName: entity.ExchangePayrollEvent,
		QueueName:    entity.QueuePayrollCalculation,
		RoutingKey:   entity.RoutingKeyPayrollCalculate,
	},
}

func (s *subscriber) setupInfra() error {
	for _, exchange := range MQExchanges {
		err := s.mq.CreateExchange(
			exchange.Name,
			exchange.Type,
			true,  // durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return err
		}
	}

	for _, queue := range MQQueue {
		_, err := s.mq.CreateQueue(
			queue.Name,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return err
		}
	}

	for _, bind := range MQQueueBind {
		err := s.mq.BindQueue(
			bind.QueueName,
			bind.ExchangeName,
			bind.RoutingKey,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *subscriber) assignEvent() {
	s.assignEventHandler(pubsub.GetEvent(entity.ExchangePayrollEvent, entity.RoutingKeyPayrollCalculate), s.uc.AttendancePeriod.PubSubGeneratePayroll)
}

func (s *subscriber) assignEventHandler(event string, handler handlerFunc) {
	ctx := context.Background()
	_, ok := s.eventHandlerMap[event]
	if ok {
		s.log.Fatal(ctx, fmt.Sprintf("Failed assign handler for event %v", event))
	}
	s.eventHandlerMap[event] = handler
}
