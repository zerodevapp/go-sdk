package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/zerodevapp/sdk-go/cmd/constants"
	useropbuilder "github.com/zerodevapp/sdk-go/cmd/useropbuilder"
)

func run4337Example() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read required environment variables
	projectID := os.Getenv("ZERODEV_PROJECT_ID")
	if projectID == "" {
		log.Fatal("ZERODEV_PROJECT_ID is required. Please set it in .env file or as an environment variable")
	}

	apiKey := os.Getenv("USEROP_BUILDER_API_KEY")
	if apiKey == "" {
		log.Fatal("USEROP_BUILDER_API_KEY is required. Please set it in .env file or as an environment variable")
	}
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
	// Create UserOpBuilder client
	//
	//
	client := useropbuilder.NewUserOpBuilder(projectID, baseURL, apiKey)

	// Optional
	client.InitialiseKernelClient(chainID, context.Background())

	// Further steps would go here...
}
