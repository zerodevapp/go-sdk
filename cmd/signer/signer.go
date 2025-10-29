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
// Returns signature as hex string (0x-prefixed) in R || S || V format.
func SignUserOpHash(userOpHash string, privateKey *ecdsa.PrivateKey) (string, error) {
	hashBytes := common.FromHex(userOpHash)
	if len(hashBytes) != 32 {
		return "", fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(hashBytes))
	}

	prefixedHash := accounts.TextHash(hashBytes)

	signatureBytes, err := crypto.Sign(prefixedHash, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign hash: %w", err)
	}

	v := signatureBytes[64]
	if v < 27 {
		v += 27
	}
	r := signatureBytes[:32]
	s := signatureBytes[32:64]

	useropSignature := append(r, s...)
	useropSignature = append(useropSignature, v)

	return "0x" + hex.EncodeToString(useropSignature), nil
}

// VerifyUserOpSignature verifies that a signature is valid for a given user operation hash.
// Returns true if the signature matches the expected address.
func VerifyUserOpSignature(userOpHash, signature, address string) (bool, error) {
	hashBytes := common.FromHex(userOpHash)
	if len(hashBytes) != 32 {
		return false, fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(hashBytes))
	}

	sigBytes := common.FromHex(signature)
	if len(sigBytes) != 65 {
		return false, fmt.Errorf("invalid signature length: expected 65 bytes, got %d", len(sigBytes))
	}

	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	digest := accounts.TextHash(hashBytes)

	pubKey, err := crypto.SigToPub(digest, sigBytes)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	expectedAddr := common.HexToAddress(address)

	return recoveredAddr == expectedAddr, nil
}

// SignAuthorization signs an EIP-7702 authorization tuple.
// EIP-7702 allows EOAs to delegate execution to a contract implementation.
// Returns a SignedAuthorization with signature components.
func SignAuthorization(chainID uint64, delegateAddressHex string, nonce uint64, privateKey *ecdsa.PrivateKey) (*types.SignedAuthorization, error) {
	addr := common.HexToAddress(delegateAddressHex)

	authTuple := []interface{}{
		chainID,
		addr,
		nonce,
	}

	rlpEncoded, err := rlp.EncodeToBytes(authTuple)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode authorization tuple: %w", err)
	}

	magic := byte(0x05)
	authMessage := append([]byte{magic}, rlpEncoded...)

	authHash := crypto.Keccak256Hash(authMessage)

	signature, err := crypto.Sign(authHash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign authorization: %w", err)
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	v := signature[64]

	yParity := int(v)
	if yParity >= 27 {
		yParity -= 27
	}

	return &types.SignedAuthorization{
		ChainID: chainID,
		Address: delegateAddressHex,
		Nonce:   nonce,
		R:       "0x" + hex.EncodeToString(r.Bytes()),
		S:       "0x" + hex.EncodeToString(s.Bytes()),
		V:       fmt.Sprintf("%d", v),
		YParity: uint8(yParity),
	}, nil
}
