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

2. Edit `.env` and add your ZeroDev project id:
   ```bash
   ZERODEV_PROJECT_ID=your-project-id-here

   # THE API KEY IS A SECRET BETWEEN THE GOLANG SERVER AND THE USEROP BUILDER SERVICE
   USEROP_BUILDER_API_KEY=your-api-key-here
   ```

Get your Project ID from the [ZeroDev Dashboard](https://dashboard.zerodev.app).


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

- An active ZeroDev project ID (from [ZeroDev Dashboard](https://dashboard.zerodev.app))
- Configure gas sponsorship policy for the ZeroDev project to allow gasless transactions.

## License

See [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- GitHub Issues: [github.com/zerodevapp/sdk-go/issues](https://github.com/zerodevapp/sdk-go/issues)
- Documentation: [docs.zerodev.app](https://docs.zerodev.app)
