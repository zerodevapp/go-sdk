// Package types defines data structures for user operations, authorizations, and API requests/responses.
package types

// Authorization represents an EIP-7702 authorization.
type Authorization struct {
	ChainID uint64 `json:"chainId"`
	Address string `json:"address"`
	Nonce   uint64 `json:"nonce"`
}

// SignedAuthorization represents an EIP-7702 authorization with signature components.
type SignedAuthorization struct {
	ChainID uint64 `json:"chainId"`
	Address string `json:"address"`
	Nonce   uint64 `json:"nonce"`
	V       string `json:"v,omitempty"` // V is optional and comes as string from API (bigint serialization)
	R       string `json:"r"`
	S       string `json:"s"`
	YParity uint8  `json:"yParity"`
}

// Call represents a single call in a user operation.
type Call struct {
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

// BuildUserOpRequest represents a request to build a user operation.
type BuildUserOpRequest struct {
	Account          string               `json:"account"`
	Authorization    *SignedAuthorization `json:"authorization,omitempty"`
	IsEip7702Account bool                 `json:"isEip7702Account,omitempty"`
	Nonce            string               `json:"nonce,omitempty"`
	Entrypoint       string               `json:"entrypoint"`
	KernelVersion    string               `json:"kernelVersion"`
	Calls            []Call               `json:"calls"`
}

// BuildUserOpResponse represents the response from building a user operation.
type BuildUserOpResponse struct {
	Sender                        string               `json:"sender"`
	Nonce                         string               `json:"nonce"`
	CallData                      string               `json:"callData"`
	AccountGasLimits              string               `json:"accountGasLimits"`
	PreVerificationGas            string               `json:"preVerificationGas"`
	GasFees                       string               `json:"gasFees"`
	PaymasterAndData              string               `json:"paymasterAndData"`
	Signature                     string               `json:"signature"`
	Factory                       string               `json:"factory,omitempty"`
	FactoryData                   string               `json:"factoryData,omitempty"`
	UserOpHash                    string               `json:"userOpHash"`
	Authorization                 *SignedAuthorization `json:"authorization,omitempty"`
	CallGasLimit                  string               `json:"callGasLimit,omitempty"`
	VerificationGasLimit          string               `json:"verificationGasLimit,omitempty"`
	MaxFeePerGas                  string               `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas          string               `json:"maxPriorityFeePerGas,omitempty"`
	Paymaster                     string               `json:"paymaster,omitempty"`
	PaymasterVerificationGasLimit string               `json:"paymasterVerificationGasLimit,omitempty"`
	PaymasterPostOpGasLimit       string               `json:"paymasterPostOpGasLimit,omitempty"`
	PaymasterData                 string               `json:"paymasterData,omitempty"`
}

// SendUserOpRequest represents a request to send a user operation.
type SendUserOpRequest struct {
	BuildUserOpResponse
	EntryPointVersion string `json:"entryPointVersion"`
	Signature         string `json:"signature"`
}

// SendUserOpResponse represents the response from sending a user operation.
type SendUserOpResponse struct {
	UserOpHash string `json:"userOpHash"`
}

// GetUserOpReceiptRequest represents a request to get a user operation receipt.
type GetUserOpReceiptRequest struct {
	UserOpHash string `json:"userOpHash"`
}

// Log represents a log entry in a transaction receipt.
// Maps to viem's Log type. Numeric fields use string to handle both hex strings and numbers.
type Log struct {
	Address          string   `json:"address"`          // The address from which this log originated
	BlockHash        string   `json:"blockHash"`        // Hash of block containing this log
	BlockNumber      string   `json:"blockNumber"`      // Number of block containing this log
	Data             string   `json:"data"`             // Contains the non-indexed arguments of the log
	LogIndex         int      `json:"logIndex"`         // Index of this log within its block
	Removed          bool     `json:"removed"`          // True if this filter has been destroyed and is invalid
	Topics           []string `json:"topics,omitempty"` // List of 0 to 4 indexed log arguments (topics)
	TransactionHash  string   `json:"transactionHash"`  // Hash of the transaction that created this log
	TransactionIndex int      `json:"transactionIndex"` // Index of the transaction that created this log
}

// TransactionReceipt represents the transaction receipt for a user operation execution.
// Maps to viem's TransactionReceipt type. Numeric fields use string to handle both formats.
type TransactionReceipt struct {
	BlobGasPrice      string `json:"blobGasPrice,omitempty"` // The actual value per gas deducted from the sender's account for blob gas (EIP-4844)
	BlobGasUsed       string `json:"blobGasUsed,omitempty"`  // The amount of blob gas used (EIP-4844)
	BlockHash         string `json:"blockHash"`              // Hash of block containing this transaction
	BlockNumber       string `json:"blockNumber"`            // Number of block containing this transaction
	ContractAddress   string `json:"contractAddress"`        // Address of new contract or null if no contract was created
	CumulativeGasUsed string `json:"cumulativeGasUsed"`      // Gas used by this and all preceding transactions in this block
	EffectiveGasPrice string `json:"effectiveGasPrice"`      // Pre-London: transaction's gasPrice. Post-London: actual gas price paid for inclusion
	From              string `json:"from"`                   // Transaction sender
	GasUsed           string `json:"gasUsed"`                // Gas used by this transaction
	Logs              []Log  `json:"logs"`                   // List of log objects generated by this transaction
	LogsBloom         string `json:"logsBloom"`              // Logs bloom filter
	Root              string `json:"root,omitempty"`         // The post-transaction state root (only for pre-Byzantium transactions)
	Status            string `json:"status"`                 // "0x1" if transaction was successful, "0x0" if it failed
	To                string `json:"to"`                     // Transaction recipient or null if deploying a contract
	TransactionHash   string `json:"transactionHash"`        // Hash of this transaction
	TransactionIndex  int    `json:"transactionIndex"`       // Index of this transaction in the block
	Type              string `json:"type"`                   // Transaction type (e.g., "0x0" for legacy, "0x2" for EIP-1559, null if not typed)
}

// UserOpReceipt represents a user operation receipt.
// Maps to viem's UserOperationReceipt type with all proper field ordering and types.
type UserOpReceipt struct {
	ActualGasCost string             `json:"actualGasCost"`       // Actual gas cost (uint256 as hex string or number)
	ActualGasUsed string             `json:"actualGasUsed"`       // Actual gas used (uint256 as hex string or number)
	EntryPoint    string             `json:"entryPoint"`          // Entrypoint address
	Logs          []Log              `json:"logs"`                // Logs emitted during execution
	Nonce         string             `json:"nonce"`               // Anti-replay parameter (nonce as uint256 hex string or number)
	Paymaster     string             `json:"paymaster,omitempty"` // Paymaster for the user operation (optional)
	Reason        string             `json:"reason,omitempty"`    // Revert reason, if unsuccessful (optional)
	Receipt       TransactionReceipt `json:"receipt"`             // Transaction receipt of the user operation execution
	Sender        string             `json:"sender"`              // Sender address
	Success       bool               `json:"success"`             // If the user operation execution was successful
	UserOpHash    string             `json:"userOpHash"`          // Hash of the user operation
}
