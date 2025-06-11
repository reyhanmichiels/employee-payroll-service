package publisher

import (
	"context"

	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/pubsub/rabbitmq"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/reyhanmichies/employee-payroll-service/src/handler/pubsub"
)

type Interface interface {
	Publish(ctx context.Context, exchangeName string, routingKey string, body interface{}) error
}

type publisher struct {
	mq   rabbitmq.Interface
	json parser.JSONInterface
}

type InitParam struct {
	MQ   rabbitmq.Interface
	Json parser.JSONInterface
}

func Init(params InitParam) Interface {
	return &publisher{
		mq:   params.MQ,
		json: params.Json,
	}
}

func (p *publisher) Publish(ctx context.Context, exchangeName string, routingKey string, body interface{}) error {
	var bodyBytes []byte
	var err error

	switch v := body.(type) {
	case string:
		bodyBytes = []byte(v)
	case []byte:
		bodyBytes = v
	default:
		bodyBytes, err = p.json.Marshal(body)
		if err != nil {
			return errors.NewWithCode(codes.CodeJSONMarshalError, "failed to marshal body: %v", err.Error())
		}
	}

	pubSubMessage := entity.PubSubMessage{
		Event:     pubsub.GetEvent(exchangeName, routingKey),
		RequestID: appcontext.GetRequestId(ctx),
		Payload:   string(bodyBytes),
	}

	messageBytes, err := p.json.Marshal(pubSubMessage)
	if err != nil {
		return errors.NewWithCode(codes.CodeJSONMarshalError, "failed to marshal message: %v", err.Error())
	}

	err = p.mq.Publish(ctx, exchangeName, routingKey, string(messageBytes))
	if err != nil {
		return err
	}

	return nil
}
