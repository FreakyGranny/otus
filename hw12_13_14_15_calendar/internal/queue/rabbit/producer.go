package rabbit

import (
	"context"
	"fmt"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/cenkalti/backoff/v3"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// Producer queue message writer.
type Producer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	reConn       chan *amqp.Error
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
}

// NewProducer returns new producer instance.
func NewProducer(uri, exchangeName, exchangeType, queue, bindingKey string) *Producer {
	return &Producer{
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queue:        queue,
		bindingKey:   bindingKey,
		reConn:       make(chan *amqp.Error),
	}
}

// Start message producer.
func (p *Producer) Start() error {
	var err error
	if err = p.connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	if err = p.channel.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}
	_, err = p.channel.QueueDeclare(
		p.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue: %s", err)
	}
	if err = p.channel.QueueBind(
		p.queue,
		p.bindingKey,
		p.exchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue: %s", err)
	}

	return nil
}

// Publish puts message to queue.
func (p *Producer) Publish(ctx context.Context, m queue.Message) error {
	select {
	case <-p.reConn:
		err := p.reConnect(ctx)
		if err != nil {
			return fmt.Errorf("reconnecting Error: %s", err)
		}
	default:
	}

	return p.channel.Publish(
		p.exchangeName,
		p.bindingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     m.ContentType,
			ContentEncoding: "",
			Body:            m.Body,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
}

func (p *Producer) connect() error {
	var err error

	p.conn, err = amqp.Dial(p.uri)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}
	p.conn.NotifyClose(p.reConn)

	return nil
}

func (p *Producer) reConnect(ctx context.Context) error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}
		time.Sleep(d)
		if err := p.connect(); err != nil {
			log.Printf("could not connect in reconnect call: %+v", err)
			continue
		}

		return nil
	}
}

// Stop stops producer.
func (p *Producer) Stop() {
	p.channel.Close()
	p.conn.Close()
}
