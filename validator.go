package erc1271

import (
	"bytes"
	"context"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/holyheld/gaelogrus"
)

type Validator struct {
	client           bind.ContractCaller
	validatorAddress common.Address
	sig              []byte
}

func NewERC1271Validator(client bind.ContractCaller) *Validator {
	return &Validator{
		client: client,
		sig:    validSignature,
	}
}

func (v *Validator) WithCustomValidSignatureHex(signature string) *Validator {
	v.sig = common.FromHex(signature)
	return v
}

func (v *Validator) WithCustomValidSignature(signature []byte) *Validator {
	v.sig = signature
	return v
}

func (v *Validator) WithValidatorAddressHex(address string) *Validator {
	v.validatorAddress = common.HexToAddress(address)
	return v
}

func (v *Validator) WithValidatorAddress(address common.Address) *Validator {
	v.validatorAddress = address
	return v
}

func (v *Validator) Validate(ctx context.Context, message []byte, signer string, signature string) (bool, error) {
	logger := gaelogrus.GetLogger(ctx).WithField("func", "Validate")
	validatorAddress := common.HexToAddress(signer)
	if len(v.validatorAddress.Bytes()) > 0 {
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
