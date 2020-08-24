package rabbitservice

// docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management
// https://github.com/rabbitmq/rabbitmq-consistent-hash-exchange
// rabbitmq-plugins enable rabbitmq_consistent_hash_exchange

import (
	"context"
	"encoding/json"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/usecases"
)

const (
	ErrManagerCreate = "can't create rabbit queue manager"
	ErrInitProducer  = "can't init producer"
	ErrInitConsumer  = "can't init consumer"
	ErrPush          = "pushing error"
	ErrPull          = "pulling error"
)

var _ entities.NotifyQueue = (*RabbitManager)(nil)

type RabbitManager struct {
	conn                    *rabbitmq.Connection
	consumeCh               *rabbitmq.Channel
	produceCh               *rabbitmq.Channel
	exchangeName, queueName string
	logger                  usecases.Logger
}

func NewRabbitManager(addr, exchangeName, queueName string, logger usecases.Logger) (*RabbitManager, error) {
	conn, err := rabbitmq.Dial(addr)
	if err != nil {
		return nil, errors.Wrap(err, ErrManagerCreate)
	}

	return &RabbitManager{
		conn:         conn,
		exchangeName: exchangeName,
		queueName:    queueName,
		logger:       logger,
	}, nil
}

func (r *RabbitManager) initProducer() error {
	sendCh, err := r.conn.Channel()
	if err != nil {
		return errors.Wrap(err, ErrInitProducer)
	}

	err = sendCh.ExchangeDeclare(r.exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, ErrInitProducer)
	}

	_, err = sendCh.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, ErrInitProducer)
	}

	if err := sendCh.QueueBind(r.queueName, "", r.exchangeName, false, nil); err != nil {
		return errors.Wrap(err, ErrInitProducer)
	}
	r.produceCh = sendCh
	return nil
}

func (r *RabbitManager) initConsumer() error {
	consumeCh, err := r.conn.Channel()
	if err != nil {
		return errors.Wrap(err, ErrInitConsumer)
	}
	r.consumeCh = consumeCh
	return nil
}

func (r *RabbitManager) Pull(ctx context.Context, alert chan<- entities.Notify) error {
	if r.produceCh == nil {
		if err := r.initConsumer(); err != nil {
			return errors.Wrap(err, ErrPull)
		}
	}
	var n entities.Notify
	var msg amqp.Delivery
	d, err := r.consumeCh.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		r.logger.Error(ctx, errors.Wrap(err, ErrPull))
	}

	loop := true
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			break
		case msg = <-d:
			err := json.Unmarshal(msg.Body, &n)
			if err != nil {
				r.logger.Error(ctx, errors.Wrap(err, ErrPull))
			}
			alert <- n
			err = msg.Ack(true)
			if err != nil {
				r.logger.Error(ctx, errors.Wrap(err, ErrPull))
			}
		}
	}
	return nil
}

func (r *RabbitManager) Push(n entities.Notify) error {
	if r.produceCh == nil {
		if err := r.initProducer(); err != nil {
			return errors.Wrap(err, ErrPush)
		}
	}
	msg, err := json.Marshal(n)
	if err != nil {
		return errors.Wrap(err, ErrPush)
	}
	err = r.produceCh.Publish(r.exchangeName, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
	if err != nil {
		return errors.Wrap(err, ErrPush)
	}
	return nil
}
