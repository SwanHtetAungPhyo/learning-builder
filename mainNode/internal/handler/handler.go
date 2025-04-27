package handler

import (
	"context"
	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/SwanHtetAungPhyo/learning/common/proto"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/avl"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/model"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"time"
)

type Handler interface {
	SubmitTx(c *fiber.Ctx) error
	GetTx(c *fiber.Ctx) error
	BlockAddition(c *fiber.Ctx) error
}

type Impl struct {
	log            *logrus.Logger
	rabbitMqClient *common.RabbitMQClient
	Validators     []model.Validator
	Consensus      *model.ConsensusResult
	chain          *common.BlockChain
}

func NewImpl(log *logrus.Logger, validators []model.Validator, chain *common.BlockChain) *Impl {
	impl := &Impl{
		log:        log,
		Validators: validators,
		chain:      chain,
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
		CreateQueue("validator1")
	impl.
		rabbitMqClient.
		BindQueueToExchange(
			selectedValidator,
			"transactions",
			"validator1key")
	return impl
}

func (i Impl) TreeBuilding() *avl.Node {
	headNode := avl.NewNode(i.Validators[0])
	if len(i.Validators) == 1 {
		return headNode
	} else {
		for _, validator := range i.Validators {
			headNode.Insert(validator)
		}
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

	//i.log.Debugln("Sending to rabbitmq: ", string(jsonBytes))
	return c.JSON(fiber.Map{
		"message":          "Transaction submitted successfully",
		"ValidatorsNumber": rand.Intn(10),
	})
}

func (i Impl) BlockAddition(c *fiber.Ctx) error {
	var incomingBlock common.Block
	if err := c.BodyParser(&incomingBlock); err != nil {
		return c.JSON(fiber.Map{
			"error":   err.Error(),
			"message": "invalid Block format",
		})
	}
	if !incomingBlock.VerifyBlockByMerkle() {
		return c.JSON(fiber.Map{
			"error": "Block is not valid",
		})
	}
	i.log.Info("Received from client: ", string(c.Body()))
	return c.JSON(fiber.Map{
		"message": "Block added successfully",
	})
}
func (i Impl) GetTx(c *fiber.Ctx) error {
	clientGrpc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer func(clientGrpc *grpc.ClientConn) {
		err := clientGrpc.Close()
		if err != nil {
			logrus.Errorln("Failed to close client connection")
		}
	}(clientGrpc)

	proto.NewBlockchainServiceClient(clientGrpc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ChainData, err := proto.NewBlockchainServiceClient(clientGrpc).GetFullChainState(ctx, &proto.Empty{})
	return c.JSON(fiber.Map{
		"message":   "Block added successfully",
		"ChainData": ChainData,
	})
}
