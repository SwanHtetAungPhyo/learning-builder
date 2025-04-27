package handler

import (
	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/avl"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/model"
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
	Validators     []model.Validator
	Consensus      *model.ConsensusResult
}

func NewImpl(log *logrus.Logger, validators []model.Validator) *Impl {
	impl := &Impl{
		log:        log,
		Validators: validators,
	}
	impl.
		Consensus = impl.
		TreeBuilding().
		CheckConsensus()
	impl.
		rabbitMqClient = common.
		NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	selectedValidator := impl.
		Consensus.
		Validators[0].Name
	impl.
		rabbitMqClient.
		Connect()
	impl.
		rabbitMqClient.
		CreateChannel()
	impl.
		rabbitMqClient.
		CreateQueue(selectedValidator)
	impl.
		rabbitMqClient.
		BindQueueToExchange(
			selectedValidator,
			"transactions",
			common.RoutingKeyCalculator(selectedValidator))
	return impl
}

func (i Impl) TreeBuilding() *avl.Node {
	headNode := avl.NewNode(i.Validators[0])
	for _, validator := range i.Validators {
		headNode.Insert(validator)
	}
	return headNode
}

func (i Impl) SubmitTx(c *fiber.Ctx) error {
	var incomingTx common.Tx

	if err := c.BodyParser(&incomingTx); err != nil {
		return c.JSON(fiber.Map{
			"error":   err.Error(),
			"message": "invalid Transaction format",
		})
	}
	i.log.Info("Received from client: ", string(c.Body()))
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
