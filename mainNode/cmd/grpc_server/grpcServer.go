package grpc_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SwanHtetAungPhyo/learning/common"
	"github.com/SwanHtetAungPhyo/learning/common/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type GrpcServer struct {
	proto.UnimplementedBlockchainServiceServer
	logger     *logrus.Logger
	blockChain *common.BlockChain
}

func NewGrpcServer(logger *logrus.Logger) *GrpcServer {
	grpcServer := &GrpcServer{logger: logger}
	grpcServer.blockChain = common.NewBlockChain("Swan")
	return grpcServer
}

func (g GrpcServer) ProposeBlockCall(ctx context.Context, request *proto.ProposeBlockRequest) (*proto.ProposeBlockResponse, error) {
	if request.Block == nil {
		return &proto.ProposeBlockResponse{}, nil
	}
	g.logger.Infoln(request.Block)
	newBlock, err := g.ConvertProtoToBlock(request.Block)
	if err != nil {
		return &proto.ProposeBlockResponse{
			ProposedBlock:  request.Block,
			ProposalStatus: "Rejected",
		}, nil
	}
	if !newBlock.VerifyBlockByMerkle() {
		return &proto.ProposeBlockResponse{
			ProposedBlock:  request.Block,
			ProposalStatus: "rejected",
		}, errors.New("signature verification failed for validation signature")
	}
	_, err = g.blockChain.AddBlock(newBlock)
	if err != nil {
		return nil, err
	}
	return &proto.ProposeBlockResponse{
		ProposedBlock:  request.Block,
		ProposalStatus: "Accepted",
	}, nil
}

func (g *GrpcServer) GetChainMetadata(ctx context.Context, _ *proto.Empty) (*proto.ChainMetadata, error) {
	g.blockChain.Mu.RLock()
	defer g.blockChain.Mu.RUnlock()

	return &proto.ChainMetadata{
		Name:       g.blockChain.ChainMetaData.Name,
		StartedAt:  g.blockChain.ChainMetaData.StartedAt,
		BlockCount: int32(len(g.blockChain.GetAllBlocks())),
		LatestHash: g.blockChain.GetLatestHash(),
	}, nil
}

func (g *GrpcServer) GetBlockByHash(ctx context.Context, req *proto.TransactionQuery) (*proto.Block, error) {
	block := g.blockChain.GetBlockByHash(req.Hash)
	if block == nil {
		return nil, status.Errorf(codes.NotFound, "block not found")
	}
	return nil, nil
}
func (g GrpcServer) ConvertProtoToBlock(protoBlock *proto.Block) (*common.Block, error) {
	if protoBlock == nil {
		return nil, fmt.Errorf("nil proto block")
	}

	// Convert BlockHeader
	var blockHeader *common.BlockHeader
	if protoBlock.BlockHeader != nil {
		blockHeader = &common.BlockHeader{
			Index:      protoBlock.BlockHeader.Index,
			MerkleRoot: protoBlock.BlockHeader.MerkleRoot,
			Validator:  protoBlock.BlockHeader.Validator,
			TimeStamp:  protoBlock.BlockHeader.Timestamp,
		}
	}

	// Convert Transactions
	txs := make([]*common.Tx, len(protoBlock.Txs))
	for i, protoTx := range protoBlock.Txs {
		if protoTx == nil {
			return nil, fmt.Errorf("nil transaction at index %d", i)
		}

		txs[i] = &common.Tx{
			From:      protoTx.From,
			To:        protoTx.To,
			Signature: protoTx.Signature,
			Amount:    int(protoTx.Amount),
			Timestamp: protoTx.Timestamp,
			PrevHash:  protoTx.PrevHash,
			Hash:      protoTx.Hash,
		}
	}

	return &common.Block{
		BlockHeader:        blockHeader,
		Hash:               protoBlock.Hash,
		PrevHash:           protoBlock.PrevHash,
		ValidatorSignature: protoBlock.ValidatorSignature,
		Txs:                txs,
	}, nil
}
func (s *GrpcServer) GetFullChainState(ctx context.Context, _ *proto.Empty) (*proto.RawChainState, error) {
	s.blockChain.Mu.RLock()
	defer s.blockChain.Mu.RUnlock()

	type ChainState struct {
		ChainMetaData struct {
			Name      string `json:"name"`
			StartedAt string `json:"startedAt"`
		} `json:"chainMetaData"`
		Blocks []*common.Block `json:"blocks"`
	}

	state := ChainState{
		ChainMetaData: s.blockChain.ChainMetaData,
		Blocks:        make([]*common.Block, len(s.blockChain.Blocks)),
	}

	copy(state.Blocks, s.blockChain.Blocks)

	jsonData, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal chain state: %v", err)
	}

	return &proto.RawChainState{
		JsonData: string(jsonData),
	}, nil
}
func (g GrpcServer) Start() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		g.logger.Fatalln(err.Error())
	}

	server := grpc.NewServer()
	proto.RegisterBlockchainServiceServer(server, &g)
	g.logger.Infoln("Grpc Server started")
	if err := server.Serve(listener); err != nil {
		g.logger.Fatalln(err.Error())
	}
}
