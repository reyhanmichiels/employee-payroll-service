package pubsub

import "fmt"

func GetEvent(exchangeName string, routingKey string) string {
	return fmt.Sprintf("%s:%s", exchangeName, routingKey)
}
