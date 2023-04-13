package erc1271

import (
	"bytes"
	"context"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/holyheld/gaelogrus"
)

// Validator is a helper struct that provides with convenience method and ERC1271-compliant validate function
type Validator struct {
	client           bind.ContractCaller
	validatorAddress common.Address
	sig              []byte
}

// NewValidator creates a new Validator instance
func NewValidator(client bind.ContractCaller) *Validator {
	return &Validator{
		client: client,
		sig:    ValidSignature,
	}
}

// WithCustomValidSignatureHex sets custom valid signature (magic value to compare the results against) using hex (string) value
func (v *Validator) WithCustomValidSignatureHex(signature string) *Validator {
	v.sig = common.FromHex(signature)
	return v
}

// WithCustomValidSignature sets custom valid signature (magic value to compare the results against) using byte slice value
func (v *Validator) WithCustomValidSignature(signature []byte) *Validator {
	v.sig = signature
	return v
}

// WithValidatorAddressHex sets validator address (target contract validator address) using hex (string) value
func (v *Validator) WithValidatorAddressHex(address string) *Validator {
	v.validatorAddress = common.HexToAddress(address)
	return v
}

// WithValidatorAddress sets validator address (target contract validator address) using common.Address value
func (v *Validator) WithValidatorAddress(address common.Address) *Validator {
	v.validatorAddress = address
	return v
}

// Validate performs all the necessary checks to tell if the signature is valid from ERC1271 standpoint
//
// Handles obvious contract (response) related errors internally, error value should be used to check if the RPC
// connection is established properly
func (v *Validator) Validate(ctx context.Context, message []byte, signer string, signature string) (bool, error) {
	logger := gaelogrus.GetLogger(ctx).WithField("func", "Validate")
	validatorAddress := common.HexToAddress(signer)
	if !IsZeroAddress(v.validatorAddress) {
		validatorAddress = v.validatorAddress
	}
	code, err := v.client.CodeAt(ctx, validatorAddress, nil)
	if err != nil {
		logger.WithField("address", validatorAddress).WithError(err).Debug("failed to check code at address")
		return false, err
	}

	isContract := len(code) > 0
	if !isContract {
		logger.WithField("address", validatorAddress).Debug("specified address is not a contract")
		return false, nil
	}

	caller, err := NewContractCaller(validatorAddress, v.client)
	if err != nil {
		logger.WithError(err).Debug("failed to create contract caller")
		return false, err
	}

	hash := common.BytesToHash(accounts.TextHash(message))
	signatureHash := common.FromHex(signature)
	res, err := caller.IsValidSignature(&bind.CallOpts{From: common.HexToAddress(signer)}, hash, signatureHash)
	if err != nil {
		logger.WithError(err).Debug("invalid signature")
		return false, nil
	}

	return bytes.Equal(res[:], v.sig), nil
}
