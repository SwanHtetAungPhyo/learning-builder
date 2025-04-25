package common

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type RabbitMQClient struct {
	Url      string
	conn     *amqp.Connection
	ch       *amqp.Channel
	connLock sync.Mutex
}

func NewRabbitMQClient(url string) *RabbitMQClient {
	return &RabbitMQClient{Url: url}
}

func (r *RabbitMQClient) Connect() *RabbitMQClient {
	r.connLock.Lock()
	defer r.connLock.Unlock()

	if r.conn == nil || r.conn.IsClosed() {
		var err error
		r.conn, err = amqp.Dial(r.Url)
		if err != nil {
			log.Panicf("Failed to connect to RabbitMQ: %v", err)
		}
		go func() {
			<-r.conn.NotifyClose(make(chan *amqp.Error))
			log.Println("RabbitMQ connection closed unexpectedly")
		}()
	}
	return r
}

func (r *RabbitMQClient) CreateChannel() *RabbitMQClient {
	r.connLock.Lock()
	defer r.connLock.Unlock()

	if r.ch == nil {
		var err error
		r.ch, err = r.conn.Channel()
		if err != nil {
			log.Panicf("Failed to open channel: %v", err)
		}
		go func() {
			<-r.ch.NotifyClose(make(chan *amqp.Error))
			log.Println("RabbitMQ channel closed unexpectedly")
		}()
	}
	return r
}

func (r *RabbitMQClient) Close() {
	r.connLock.Lock()
	defer r.connLock.Unlock()

	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
		r.ch = nil
	}

	if r.conn != nil && !r.conn.IsClosed() {
		if err := r.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
		r.conn = nil
	}
}

func (r *RabbitMQClient) CreateQueue(queueName string) {
	if r.ch == nil {
		r.CreateChannel()
	}
	_, err := r.ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to declare queue: %v", err)
	}
}

func (r *RabbitMQClient) SendMsgText(message string, key string) {
	if r.ch == nil {
		r.CreateChannel()
	}
	err := r.ch.Publish("", key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	if err != nil {
		log.Panicf("Failed to publish message: %v", err)
	}
}

func (r *RabbitMQClient) SendMsgJson(message []byte, exchangeName, key string) {
	if r.ch == nil {
		r.CreateChannel()
	}
	err := r.ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to declare exchange: %v", err)
	}
	err = r.ch.Publish(exchangeName, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	if err != nil {
		log.Panicf("Failed to publish message: %v", err)
	}
}

func (r *RabbitMQClient) ConsumeMsgWithKey(queueName, exchange, routingKey string) <-chan amqp.Delivery {
	if r.ch == nil {
		r.CreateChannel()
	}
	msgs, err := r.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to consume messages: %v", err)
	}
	return msgs
}

func (r *RabbitMQClient) ConsumeMsgWithExchange(queueName, exchangeName, exchangeType, routingKey string) <-chan amqp.Delivery {
	if r.ch == nil {
		r.CreateChannel()
	}
	err := r.ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to declare exchange: %v", err)
	}
	_, err = r.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to declare queue: %v", err)
	}
	err = r.ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Panicf("Failed to bind queue: %v", err)
	}
	msgs, err := r.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to consume messages: %v", err)
	}
	return msgs
}

func (r *RabbitMQClient) CreateExchange(name string, s string) {
	if r.ch == nil {
		r.CreateChannel()
	}
	err := r.ch.ExchangeDeclare(name, s, true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to declare exchange: %v", err)
	}
}

func (r *RabbitMQClient) BindQueueToExchange(queueName, exchangeName, routingKey string) {
	if r.ch == nil {
		r.CreateChannel()
	}
	err := r.ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Panicf("Failed to bind queue: %v", err)
	}
}
