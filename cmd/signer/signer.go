package signer

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zerodevapp/sdk-go/cmd/types"
)

// SignUserOpHash signs a user operation hash using Ethereum's personal_sign format.
// This uses go-ethereum's accounts.TextHash for standard Ethereum message signing.
// Format: sign(keccak256("\x19Ethereum Signed Message:\n32" + hash))
//
// Parameters:
//   - userOpHash: The user operation hash (with or without 0x prefix)
//   - privateKey: The ECDSA private key to sign with
//
// Returns:
//   - The signature as a hex string (with 0x prefix) in the format R || S || V
//   - An error if signing fails
func SignUserOpHash(userOpHash string, privateKey *ecdsa.PrivateKey) (string, error) {
	// Parse hash using go-ethereum's common.FromHex (handles 0x prefix automatically)
	hashBytes := common.FromHex(userOpHash)
	if len(hashBytes) != 32 {
		return "", fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(hashBytes))
	}

	// Use go-ethereum's standard TextHash for personal_sign format
	// This applies the "\x19Ethereum Signed Message:\n" prefix
	prefixedHash := accounts.TextHash(hashBytes)

	// Sign using go-ethereum's crypto.Sign
	// Returns signature in [R || S || V] format where V is 0 or 1
	signatureBytes, err := crypto.Sign(prefixedHash, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign hash: %w", err)
	}

	// crypto.Sign returns [R || S || V] where V is already 0 or 1

	v := signatureBytes[64]
	if v < 27 {
		v += 27
	}
	r := signatureBytes[:32]
	s := signatureBytes[32:64]

	// Rearrange signature to encodePacked(r, s, v)
	useropSignature := append(r, s...)
	useropSignature = append(useropSignature, v)

	return "0x" + hex.EncodeToString(useropSignature), nil

}

// VerifyUserOpSignature verifies that a signature is valid for a given user operation hash.
// This uses go-ethereum's standard signature verification with Ecrecover.
//
// Parameters:
//   - userOpHash: The user operation hash (with or without 0x prefix)
//   - signature: The signature to verify (with or without 0x prefix)
//   - address: The expected signer address
//
// Returns:
//   - true if the signature is valid, false otherwise
//   - An error if verification fails
func VerifyUserOpSignature(userOpHash, signature, address string) (bool, error) {
	// Parse hash using common.FromHex (handles 0x prefix automatically)
	hashBytes := common.FromHex(userOpHash)
	if len(hashBytes) != 32 {
		return false, fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(hashBytes))
	}

	// Parse signature using common.FromHex
	sigBytes := common.FromHex(signature)
	if len(sigBytes) != 65 {
		return false, fmt.Errorf("invalid signature length: expected 65 bytes, got %d", len(sigBytes))
	}

	// Normalize V value to 0 or 1 if needed
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Use go-ethereum's standard TextHash for message prefix
	digest := accounts.TextHash(hashBytes)

	// Recover public key using go-ethereum's SigToPub
	pubKey, err := crypto.SigToPub(digest, sigBytes)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Derive address from public key using go-ethereum's PubkeyToAddress
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Parse expected address using common.HexToAddress (case-insensitive)
	expectedAddr := common.HexToAddress(address)

	return recoveredAddr == expectedAddr, nil
}

// SignAuthorization signs an EIP-7702 authorization tuple.
// EIP-7702 allows EOAs to temporarily delegate their code to a smart contract implementation.
// This uses go-ethereum's rlp.EncodeToBytes for standard RLP encoding.
//
// The authorization tuple format is: keccak256(MAGIC || rlp([chainId, address, nonce]))
// where MAGIC = 0x05
//
// Parameters:
//   - chainID: The chain ID for the authorization
//   - contractAddress: The address of the contract implementation to delegate to
//   - nonce: The nonce for the authorization (typically 0 for first authorization)
//   - privateKey: The ECDSA private key to sign with
//
// Returns:
//   - Authorization struct with R, S, V, and YParity values
//   - An error if signing fails
func SignAuthorization(chainID uint64, delegateAddressHex string, nonce uint64, privateKey *ecdsa.PrivateKey) (*types.SignedAuthorization, error) {
	// Parse address to address for encoding
	addr := common.HexToAddress(delegateAddressHex)

	// Create authorization tuple for RLP encoding
	// EIP-7702 specifies: [chainId, address, nonce]
	authTuple := []interface{}{
		chainID,
		addr,
		nonce,
	}

	// Use go-ethereum's rlp.EncodeToBytes for standard RLP encoding
	rlpEncoded, err := rlp.EncodeToBytes(authTuple)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode authorization tuple: %w", err)
	}

	// Build the authorization message according to EIP-7702
	// Format: MAGIC || rlp([chainId, address, nonce])
	// MAGIC = 0x05 for EIP-7702
	magic := byte(0x05)
	authMessage := append([]byte{magic}, rlpEncoded...)

	// Hash the authorization message using go-ethereum's Keccak256Hash
	authHash := crypto.Keccak256Hash(authMessage)

	// Sign using go-ethereum's crypto.Sign
	signature, err := crypto.Sign(authHash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign authorization: %w", err)
	}

	// Extract R, S, V from signature
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	v := signature[64]

	yParity := int(v)
	if yParity >= 27 {
		yParity -= 27
	}

	return &types.SignedAuthorization{
		ChainID: chainID,
		Address: delegateAddressHex, // Use go-ethereum's Hex() method
		Nonce:   nonce,
		R:       "0x" + hex.EncodeToString(r.Bytes()),
		S:       "0x" + hex.EncodeToString(s.Bytes()),
		V:       fmt.Sprintf("%d", v),
		YParity: uint8(yParity),
	}, nil
}
