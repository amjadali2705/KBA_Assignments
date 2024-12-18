package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func blockEventListener(organization, channelName string) {

	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	fmt.Println("****Start Block Event Listening****")

	events, err := network.BlockEvents(ctx, client.WithStartBlock(4))

	if err != nil {
		panic(fmt.Errorf("failed to start block event listening: %w", err))
	}

	for event := range events {
		fmt.Printf("Recieved block number %d\n", event.GetHeader().GetNumber())
	}
}

func chaincodeEventListener(organization, channelName, chaincodeName string) {
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	fmt.Println("**Started Chaincode Event Listener")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)

	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	for event := range events {
		fmt.Printf("Recieved Chaincode Event: %s\n Data: %s\n \n", event.EventName, event.Payload)
	}
}

func pvtBlockListener(organization, channelName string) {
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	fmt.Println("****Started Listening for Private txn in Blocks")

	events, err := network.BlockAndPrivateDataEvents(ctx, client.WithStartBlock(1))

	if err != nil {
		panic(fmt.Errorf("failed to start private data event listening: %w", err))

	}
	for event := range events {
		if event.GetPrivateDataMap() != nil {
			fmt.Printf("Received block %d with pvt data \n", event.Block.GetHeader().GetNumber())
		}
	}
}
