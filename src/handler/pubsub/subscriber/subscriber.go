package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/operator"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/pubsub/rabbitmq"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase"
)

type handlerFunc func(ctx context.Context, payload entity.PubSubMessage) error

type Interface interface {
	Subscribe()
}

type subscriber struct {
	mq              rabbitmq.Interface
	log             log.Interface
	json            parser.JSONInterface
	uc              *usecase.Usecases
	eventHandlerMap map[string]handlerFunc
}

type InitParam struct {
	MQ   rabbitmq.Interface
	Log  log.Interface
	Json parser.JSONInterface
	UC   *usecase.Usecases
}

func Init(params InitParam) Interface {
	pubSub := &subscriber{
		mq:              params.MQ,
		log:             params.Log,
		json:            params.Json,
		uc:              params.UC,
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

func (s *subscriber) Subscribe() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.log.Panic(err)
			}
		}()

		ctx := s.addFieldsToContext(context.Background(), "runner", uuid.New().String())
		for _, queue := range MQQueue {
			s.mq.Subscribe(ctx, queue.Name, s.runnerWrapper)
		}

	}()
}

func (s *subscriber) runnerWrapper(ctx context.Context, exchangeName string, routingKey string, message string) error {
	pubSubMessage, err := s.unmarshalPubSubMessage(message)
	if err != nil {
		return err
	}

	ctx = s.addFieldsToContext(ctx, pubSubMessage.Event, pubSubMessage.RequestID)

	defer func() {
		startTime := appcontext.GetRequestStartTime(ctx)
		s.log.Info(ctx, fmt.Sprintf("finish running handler for event %s done in %v", pubSubMessage.Event, time.Since(startTime)))
	}()

	s.log.Info(ctx, fmt.Sprintf("running handler for event %s", pubSubMessage.Event))

	handler, ok := s.eventHandlerMap[pubSubMessage.Event]
	if !ok {
		s.log.Warn(ctx, fmt.Sprintf("no handler found for event %s", pubSubMessage.Event))
		return nil
	}

	err = handler(ctx, pubSubMessage)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("error handling event %s: %v", pubSubMessage.Event, err))
		return err
	}

	return nil
}

func (s *subscriber) addFieldsToContext(ctx context.Context, event, requestID string) context.Context {
	ctx = appcontext.SetUserAgent(ctx, fmt.Sprintf("subscriber event %v", event))
	ctx = appcontext.SetRequestId(ctx, operator.Ternary(requestID == "", uuid.New().String(), requestID))
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())
	return ctx
}

func (s *subscriber) unmarshalPubSubMessage(message string) (entity.PubSubMessage, error) {
	var pubSubMessage entity.PubSubMessage
	err := s.json.Unmarshal([]byte(message), &pubSubMessage)
	if err != nil {
		return entity.PubSubMessage{}, errors.NewWithCode(codes.CodeJSONUnmarshalError, "failed to unmarshal pubsub message: %v", err.Error())
	}
	return pubSubMessage, nil
}
