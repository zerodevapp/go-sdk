# ZeroDev Go SDK

Go SDK for building and sending EIP-4337 User Operations with EIP-7702 support via ZeroDev's Kernel accounts.

## Installation

```bash
go get github.com/zerodevapp/sdk-go
```

## Features

- Build and send EIP-4337 User Operations
- EIP-7702 authorization signing for EOA delegation
- Support for multiple Kernel versions (0.3.1, 0.3.2, 0.3.3)
- Batch multiple calls in a single User Operation
- Wait for User Operation receipts with automatic polling
- ECDSA signature support

## Environment Variables

The SDK requires the following environment variables to be set:

1. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` and add your ZeroDev credentials:
   ```bash
   ZERODEV_PROJECT_ID=your-project-id-here
   USEROP_BUILDER_API_KEY=your-api-key-here
   ```

Get your Project ID and API Key from the [ZeroDev Dashboard](https://dashboard.zerodev.app).

Alternatively, you can set these as system environment variables:
```bash
export ZERODEV_PROJECT_ID="your-project-id"
export USEROP_BUILDER_API_KEY="your-api-key"
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/ethereum/go-ethereum/crypto"
    "github.com/joho/godotenv"
    "github.com/zerodevapp/sdk-go/cmd/constants"
    "github.com/zerodevapp/sdk-go/cmd/signer"
    "github.com/zerodevapp/sdk-go/cmd/types"
    "github.com/zerodevapp/sdk-go/cmd/useropbuilder"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Configuration
    projectID := os.Getenv("ZERODEV_PROJECT_ID")
    apiKey := os.Getenv("USEROP_BUILDER_API_KEY")
    chainID := uint64(11155111) // Sepolia testnet
    kernelVersion := constants.KernelVersion033
    baseURL := "https://api.zerodev.app"
    entrypointVersion := "0.7"

    // Generate or load your private key
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        log.Fatal(err)
    }
    address := crypto.PubkeyToAddress(privateKey.PublicKey)

    // Get account implementation address for the kernel version
    accountImplementation, err := constants.GetAccountImplementationAddress(kernelVersion)
    if err != nil {
        log.Fatal(err)
    }

    // Sign EIP-7702 authorization (only needed for first transaction)
    authorization, err := signer.SignAuthorization(chainID, accountImplementation, 0, privateKey)
    if err != nil {
        log.Fatal(err)
    }

    // Create client
    client := useropbuilder.NewUserOpBuilder(projectID, baseURL, apiKey)

    // Optional: Initialize kernel client
    client.InitialiseKernelClient(chainID, context.Background())

    // Define calls
    calls := []types.Call{
        {
            To:    "0x0000000000000000000000000000000000000000",
            Value: "0",
            Data:  "0x",
        },
    }

    // Build User Operation
    buildReq := &types.BuildUserOpRequest{
        Account:          address.Hex(),
        Authorization:    authorization,
        IsEip7702Account: true,
        Entrypoint:       entrypointVersion,
        KernelVersion:    string(kernelVersion),
        Calls:            calls,
    }

    buildResp, err := client.BuildUserOp(context.Background(), chainID, buildReq)
    if err != nil {
        log.Fatal(err)
    }

    // Sign User Operation hash
    signature, err := signer.SignUserOpHash(buildResp.UserOpHash, privateKey)
    if err != nil {
        log.Fatal(err)
    }

    // Send User Operation
    sendReq := &types.SendUserOpRequest{
        BuildUserOpResponse: *buildResp,
        EntryPointVersion:   entrypointVersion,
        Signature:           signature,
    }

    sendResp, err := client.SendUserOp(context.Background(), chainID, sendReq)
    if err != nil {
        log.Fatal(err)
    }

    // Wait for receipt
    receiptReq := &types.GetUserOpReceiptRequest{
        UserOpHash: sendResp.UserOpHash,
    }

    receipt, err := client.WaitForUserOpReceipt(
        context.Background(),
        chainID,
        receiptReq,
        2*time.Second,  // poll interval
        60*time.Second, // timeout
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Success! UserOp Hash: %s", receipt.UserOpHash)
}
```

## Running Examples

First, set up your environment variables:

```bash
# Copy the example .env file
cp .env.example .env
# Edit .env and add your ZERODEV_PROJECT_ID and USEROP_BUILDER_API_KEY
```

Then run an example:

```bash
# Run the EIP-7702 example
go run ./example -e 7702

# Show available commands
go run ./example help
```

## Requirements

- Go 1.25.3 or higher
- An active ZeroDev project ID and API key

## License

See [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- GitHub Issues: [github.com/zerodevapp/sdk-go/issues](https://github.com/zerodevapp/sdk-go/issues)
- Documentation: [docs.zerodev.app](https://docs.zerodev.app)
