package rabbit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// Consumer rmq consumer.
type Consumer struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	done        chan error
	wg          *sync.WaitGroup
	consumerTag string
	uri         string
	queue       string
}

// NewConsumer returns new consumer instance.
func NewConsumer(consumerTag, uri, queue string) *Consumer {
	return &Consumer{
		consumerTag: consumerTag,
		uri:         uri,
		queue:       queue,
		done:        make(chan error),
		wg:          &sync.WaitGroup{},
	}
}

// Worker consumer worker func.
type Worker func(context.Context, *sync.WaitGroup, <-chan amqp.Delivery)

// Handle process incoming messages.
func (c *Consumer) Handle(ctx context.Context, fn Worker, threads int) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for {
		for i := 0; i < threads; i++ {
			c.wg.Add(1)
			go fn(ctx, c.wg, msgs)
		}

		err, open := <-c.done
		c.wg.Wait()
		if !open {
			return nil
		}
		if err != nil {
			msgs, err = c.reConnect(ctx)
			if err != nil {
				return fmt.Errorf("reconnecting Error: %s", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}

func (c *Consumer) connect() error {
	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("channel Closed")
	}()

	return nil
}

// Задекларировать очередь, которую будем слушать.
func (c *Consumer) announceQueue() (<-chan amqp.Delivery, error) {
	err := c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("error setting qos: %s", err)
	}
	_, err = c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %s", err)
	}
	msgs, err := c.channel.Consume(
		c.queue,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	return msgs, nil
}

func (c *Consumer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		time.Sleep(d)
		if err := c.connect(); err != nil {
			log.Printf("could not connect in reconnect call: %+v", err)
			continue
		}

		msgs, err := c.announceQueue()
		if err != nil {
			fmt.Printf("Couldn't connect: %+v", err)
			continue
		}

		return msgs, nil
	}
}

// Stop stops consumer.
func (c *Consumer) Stop() {
	close(c.done)
}
