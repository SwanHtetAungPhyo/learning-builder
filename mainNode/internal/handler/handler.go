package handler

import (
	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type Handler interface {
	SubmitTx(c *fiber.Ctx) error
	GetTx(c *fiber.Ctx) error
}

type Impl struct {
	log            *logrus.Logger
	rabbitMqClient *common.RabbitMQClient
}

func NewImpl(log *logrus.Logger) *Impl {
	impl := &Impl{
		log: log,
	}
	impl.rabbitMqClient = common.NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	impl.rabbitMqClient.Connect()
	impl.rabbitMqClient.CreateChannel()
	impl.rabbitMqClient.CreateQueue("validator1")
	impl.rabbitMqClient.BindQueueToExchange("validator1", "transactions", "validator1key")
	return impl
}

func (i Impl) SubmitTx(c *fiber.Ctx) error {
	var incomingTx common.Tx
	i.log.Infoln("Received request: ", c.Body())
	if err := c.BodyParser(&incomingTx); err != nil {
		return c.JSON(fiber.Map{
			"error":   err.Error(),
			"message": "invalid Transaction format",
		})
	}
	jsonBytes := common.Must[[]byte](json.Marshal(incomingTx))

	i.rabbitMqClient.SendMsgJson(jsonBytes, "transactions", "validator1key")

	i.log.Debugln("Sending to rabbitmq: ", string(jsonBytes))
	return c.JSON(fiber.Map{
		"message":          "Transaction submitted successfully",
		"ValidatorsNumber": rand.Intn(10),
	})
}

func (i Impl) GetTx(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
