package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/operator"
	"github.com/reyhanmichiels/go-pkg/v2/pubsub/rabbitmq"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

type handlerFunc func(ctx context.Context, payload entity.PubSubMessage) error

type Interface interface {
	Subscribe()
}

type pubSub struct {
	mq              rabbitmq.Interface
	log             log.Interface
	eventHandlerMap map[string]handlerFunc
}

type InitParam struct {
	MQConfig rabbitmq.Config
	Log      log.Interface
}

func Init(params InitParam) Interface {
	pubSub := &pubSub{
		mq:              rabbitmq.Init(params.MQConfig, params.Log),
		log:             params.Log,
		eventHandlerMap: make(map[string]handlerFunc),
	}

	err := pubSub.setupInfra()
	if err != nil {
		pubSub.log.Fatal(context.Background(), fmt.Sprintf("[FATAL] failed to setup pubsub infrastructure: %s", err.Error()))
	}

	pubSub.assignEvent()

	go func() {
		pubSub.mq.MonitorConnection()
		pubSub.Subscribe()
	}()

	return pubSub
}

func (p *pubSub) Subscribe() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				p.log.Panic(err)
			}
		}()

		ctx := p.addFieldsToContext(context.Background(), "runner", uuid.New().String())
		for _, queue := range MQQueue {
			p.mq.Subscribe(ctx, queue.Name, p.runnerWrapper)
		}

	}()
}

func (p *pubSub) runnerWrapper(ctx context.Context, exchangeName string, routingKey string, message string) error {
	pubSubMessage, err := p.unmarshalPubSubMessage(message)
	if err != nil {
		return err
	}

	ctx = p.addFieldsToContext(ctx, pubSubMessage.Event, pubSubMessage.RequestID)

	defer func() {
		startTime := appcontext.GetRequestStartTime(ctx)
		p.log.Info(ctx, fmt.Sprintf("finish running handler for event %s done in %v", pubSubMessage.Event, time.Since(startTime)))
	}()

	p.log.Info(ctx, fmt.Sprintf("running handler for event %s", pubSubMessage.Event))

	handler, ok := p.eventHandlerMap[pubSubMessage.Event]
	if !ok {
		p.log.Warn(ctx, fmt.Sprintf("no handler found for event %s", pubSubMessage.Event))
		return nil
	}

	err = handler(ctx, pubSubMessage)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("error handling event %s: %v", pubSubMessage.Event, err))
		return err
	}

	return nil
}

func (p *pubSub) addFieldsToContext(ctx context.Context, event, requestID string) context.Context {
	ctx = appcontext.SetUserAgent(ctx, fmt.Sprintf("subscriber event %v", event))
	ctx = appcontext.SetRequestId(ctx, operator.Ternary(requestID == "", uuid.New().String(), requestID))
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())
	return ctx
}

func GetEvent(exchangeName string, routingKey string) string {
	return fmt.Sprintf("%s:%s", exchangeName, routingKey)
}

func (p *pubSub) unmarshalPubSubMessage(message string) (entity.PubSubMessage, error) {
	var pubSubMessage entity.PubSubMessage
	err := json.Unmarshal([]byte(message), &pubSubMessage)
	if err != nil {
		return entity.PubSubMessage{}, errors.NewWithCode(codes.CodeJSONUnmarshalError, "failed to unmarshal pubsub message: %v", err.Error())
	}
	return pubSubMessage, nil
}
