package pubsub

import (
	"context"
	"fmt"

	"github.com/reyhanmichiels/go-pkg/v2/pubsub/rabbitmq"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
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

func (p *pubSub) setupInfra() error {
	for _, exchange := range MQExchanges {
		err := p.mq.CreateExchange(
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
		_, err := p.mq.CreateQueue(
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
		err := p.mq.BindQueue(
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

func (p *pubSub) assignEvent() {
	p.assignEventHandler(GetEvent(entity.ExchangePayrollEvent, entity.RoutingKeyPayrollCalculate), func(ctx context.Context, payload entity.PubSubMessage) error {

		fmt.Printf("\n\n\n SELAMAT KAMU BERHASIL \n\n\n")

		return nil
	})
}

func (p *pubSub) assignEventHandler(event string, handler handlerFunc) {
	ctx := context.Background()
	_, ok := p.eventHandlerMap[event]
	if ok {
		p.log.Fatal(ctx, fmt.Sprintf("Failed assign handler for event %v", event))
	}
	p.eventHandlerMap[event] = handler
}
