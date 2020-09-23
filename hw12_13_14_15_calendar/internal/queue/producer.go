package queue

import "context"

// Producer message queue writer.
type Producer interface {
	Start() error
	Publish(context.Context, Message) error
	Stop()
}
