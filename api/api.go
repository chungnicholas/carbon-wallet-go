package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	log "github.com/sirupsen/logrus"
)

// APIs in this file is added in a "if we need it, then we add it basis"

// GetAccount gets an account from its bech32 address
func GetAccount(targetGRPCAddress string, bech32Address string, clientCtx client.Context) (account *authtypes.BaseAccount, err error) {
	grpcConn, err := GetGRPCConnection(targetGRPCAddress, clientCtx)
	if err != nil {
		return nil, err
	}
	defer grpcConn.Close()

	//log.Info("Getting account: ", bech32Address)

	// This creates a gRPC client to query the x/account service.
	authClient := authtypes.NewQueryClient(grpcConn)
	accountRes, err := authClient.Account(
		context.Background(),
		&authtypes.QueryAccountRequest{
			Address: bech32Address,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}

	ba := authtypes.BaseAccount{}
	err = ba.Unmarshal(accountRes.Account.Value)
	if err != nil {
		log.Error(err)
		return
	}

	//log.Infof("Found Account. Address: %s, AccountNumber: %d, Sequence: %d", ba.Address, ba.AccountNumber, ba.Sequence)

	return &ba, nil
}

// GetChainID of the node from tendermint grpc
func GetChainID(targetGRPCAddress string, clientCtx client.Context) (chainID string, err error) {
	grpcConn, err := GetGRPCConnection(targetGRPCAddress, clientCtx)
	if err != nil {
		return "", err
	}
	defer grpcConn.Close()
	log.Info("Getting Node Info")

	serviceClient := tmservice.NewServiceClient(grpcConn)
	nodeInfoRes, err := serviceClient.GetNodeInfo(
		context.Background(),
		&tmservice.GetNodeInfoRequest{},
	)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("node info: ", nodeInfoRes.DefaultNodeInfo)
	return nodeInfoRes.DefaultNodeInfo.Network, nil
}

// GetLatestBlockHeight of the node from tendermint grpc
func GetLatestBlockHeight(targetGRPCAddress string, clientCtx client.Context) (height int64, err error) {
	grpcConn, err := GetGRPCConnection(targetGRPCAddress, clientCtx)
	if err != nil {
		return 0, err
	}
	defer grpcConn.Close()

	serviceClient := tmservice.NewServiceClient(grpcConn)
	lastestBlockRes, err := serviceClient.GetLatestBlock(
		context.Background(),
		&tmservice.GetLatestBlockRequest{},
	)
	height = lastestBlockRes.Block.Header.Height
	if err != nil {
		log.Error(err)
		return 0, err
	}
	if height <= 0 {
		err = errors.New(fmt.Sprintf("get latest block height is invalid: %+v\n", height))
		log.Error(err)
		return 0, err
	}

	return height, nil
}
