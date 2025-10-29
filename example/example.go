package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zerodevapp/sdk-go/cmd/constants"
	"github.com/zerodevapp/sdk-go/cmd/signer"
	types "github.com/zerodevapp/sdk-go/cmd/types"
	useropbuilder "github.com/zerodevapp/sdk-go/cmd/useropbuilder"
)

func logJSON(v interface{}) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return
	}
	log.Printf("%s", string(jsonBytes))
}

func main() {
	projectID := "PROJECT_ID"
	apiKey := "your-api-key-here"
	chainID := uint64(11155111) // Sepolia
	kernelVersion := constants.KernelVersion033
	baseURL := "http://localhost:3010"
	entrypointVersion := "0.7"
	//
	//
	// Get account implementation address from SDK
	//
	//
	accountImplementationAddress, err := constants.GetAccountImplementationAddress(kernelVersion)
	if err != nil {
		log.Fatalf("Failed to get account implementation address: %v", err)
	}
	fmt.Println("=== Configuration ===")
	fmt.Printf("\tProject ID: %s\n", projectID)
	fmt.Printf("\tChain ID: %d (Sepolia)\n", chainID)
	fmt.Printf("\tKernel Version: %s\n", kernelVersion)
	fmt.Printf("\tEntrypoint Version: %s\n", entrypointVersion)
	fmt.Printf("\tAccount Implementation Address: %s\n", accountImplementationAddress)
	fmt.Printf("\tBase URL: %s\n", baseURL)

	//
	//
	// Get ECDSA private key for signing
	//
	//
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	// Ethereum address from private key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressHex := address.Hex()
	fmt.Println("\n=== Account ===")
	fmt.Println("\tAddress:", addressHex)

	//
	//
	// Sign EIP-7702 authorization
	//
	//
	authorization, err := signer.SignAuthorization(chainID, accountImplementationAddress, 0, privateKey)
	if err != nil {
		log.Fatalf("Failed to sign authorization: %v", err)
	}
	fmt.Println("\n=== Authorization Signed ===")
	fmt.Printf("\tR: %s\n", authorization.R)
	fmt.Printf("\tS: %s\n", authorization.S)
	fmt.Printf("\tV: %s\n", authorization.V)
	fmt.Printf("\tYParity: %d\n", authorization.YParity)

	//
	//
	// Create UserOpBuilder client
	//
	//
	client := useropbuilder.NewUserOpBuilder(projectID, baseURL, apiKey)

	// Optional
	client.InitialiseKernelClient(chainID, context.Background())

	//
	//
	// Define calls to be included in the user operation
	//
	//
	calls := []types.Call{
		{
			To:    "0x0000000000000000000000000000000000000000",
			Value: "0",
			Data:  "0x",
		},
		{
			To:    "0x0000000000000000000000000000000000000001",
			Value: "0",
			Data:  "0x",
		},
	}

	//
	//
	// Build user operation
	//
	//
	fmt.Println("\n\n\n=== Step 1: Build User Operation ===")
	buildReq := &types.BuildUserOpRequest{
		Account:          addressHex, // Use .Hex() to convert common.Address to string
		Authorization:    authorization,
		IsEip7702Account: true,
		Entrypoint:       entrypointVersion,
		KernelVersion:    string(kernelVersion),
		Calls:            calls,
	}
	logJSON(buildReq)

	buildUseropResponse, err := client.BuildUserOp(context.Background(), chainID, buildReq)
	if err != nil {
		log.Fatalf("Failed to build user op: %v", err)
	}

	fmt.Printf("\n✓ UserOp built successfully!\n")
	logJSON(buildUseropResponse)

	//
	//
	// Sign the user operation hash
	//
	//
	fmt.Println("\n\n\n=== Step 2: Sign User Operation ===")
	fmt.Printf("Signing hash: %s\n", buildUseropResponse.UserOpHash)

	// Sign using go-ethereum's crypto.Sign with personal_sign format
	signatureHex, err := signer.SignUserOpHash(buildUseropResponse.UserOpHash, privateKey)
	if err != nil {
		log.Fatalf("Failed to sign user op hash: %v", err)
	}
	fmt.Printf("✓ UserOp hash signed successfully!\n")

	//
	//
	// Send the built and signed user operation
	//
	//
	fmt.Println("\n\n\n=== Step 3: Send User Operation ===")

	// Send user operation
	sendUseropResponse := &types.SendUserOpRequest{
		BuildUserOpResponse: *buildUseropResponse,
		EntryPointVersion:   entrypointVersion,
		Signature:           signatureHex,
	}

	sendResp, err := client.SendUserOp(context.Background(), chainID, sendUseropResponse)
	if err != nil {
		log.Fatalf("Failed to send user op: %v", err)
	}

	fmt.Printf("\n✓ UserOp sent successfully!\n")
	logJSON(sendResp)

	//
	//
	// Wait for user operation receipt
	//
	//
	fmt.Println("\n\n\n=== Step 4: Wait for Receipt ===")
	// Wait for receipt
	receiptReq := &types.GetUserOpReceiptRequest{
		UserOpHash: sendResp.UserOpHash,
	}
	receipt, err := client.WaitForUserOpReceipt(context.Background(), chainID, receiptReq, 2*time.Second, 60*time.Second)
	if err != nil {
		log.Fatalf("Failed to get user op receipt: %v", err)
	}

	fmt.Println("\n=== Result ===")
	fmt.Printf("✓ UserOp receipt received!\n")
	logJSON(receipt)

	//
	//
	// Build request 2 - without the authorization
	//
	//
	fmt.Println("\n\n\n=== Step 5: Build User Operation Without Authorization ===")
	buildReq2 := &types.BuildUserOpRequest{
		Account:          addressHex, // Use .Hex() to convert common.Address to string
		IsEip7702Account: true,
		Entrypoint:       entrypointVersion,
		KernelVersion:    string(kernelVersion),
		Calls:            calls,
	}
	logJSON(buildReq2)

	buildUseropResponse2, err := client.BuildUserOp(context.Background(), chainID, buildReq2)
	if err != nil {
		log.Fatalf("Failed to build user op: %v", err)
	}

	fmt.Printf("\n✓ UserOp built successfully!\n")
	logJSON(buildUseropResponse2)

	//
	//
	// Sign the user operation hash
	//
	//
	fmt.Println("\n\n\n=== Step 6: Sign User Operation Without Authorization ===")
	fmt.Printf("Signing hash: %s\n", buildUseropResponse2.UserOpHash)

	// Sign using go-ethereum's crypto.Sign with personal_sign format
	signatureHex2, err := signer.SignUserOpHash(buildUseropResponse2.UserOpHash, privateKey)
	if err != nil {
		log.Fatalf("Failed to sign user op hash: %v", err)
	}
	fmt.Printf("✓ UserOp hash signed successfully!\n")

	//
	//
	// Send the built and signed user operation
	//
	//
	fmt.Println("\n\n\n=== Step 7: Send User Operation Without Authorization ===")

	// Send user operation
	sendUseropResponse2 := &types.SendUserOpRequest{
		BuildUserOpResponse: *buildUseropResponse2,
		EntryPointVersion:   entrypointVersion,
		Signature:           signatureHex2,
	}

	sendResp2, err := client.SendUserOp(context.Background(), chainID, sendUseropResponse2)
	if err != nil {
		log.Fatalf("Failed to send user op: %v", err)
	}

	fmt.Printf("\n✓ UserOp sent successfully!\n")
	logJSON(sendResp2)

	//
	//
	// Wait for user operation receipt
	//
	//
	fmt.Println("\n\n\n=== Step 8: Wait for Receipt Without Authorization ===")
	// Wait for receipt
	receiptReq2 := &types.GetUserOpReceiptRequest{
		UserOpHash: sendResp2.UserOpHash,
	}
	receipt2, err := client.WaitForUserOpReceipt(context.Background(), chainID, receiptReq2, 2*time.Second, 60*time.Second)
	if err != nil {
		log.Fatalf("Failed to get user op receipt: %v", err)
	}

	fmt.Println("\n=== Result ===")
	fmt.Printf("✓ UserOp receipt received!\n")
	logJSON(receipt2)
}
