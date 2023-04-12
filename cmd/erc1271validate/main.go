package main

import (
	"context"
	"github.com/holyheld/gaelogrus"
	"os"

	"flag"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holyheld/erc1271"
)

func main() {
	var message string
	var signer string
	var signature string
	var rpcURL string
	var validatorAddress string
	var customValidSignature string
	var debug bool

	flag.StringVar(&rpcURL, "rpc", "https://cloudflare-eth.com", "specifies rpc url explicitly")
	flag.StringVar(&rpcURL, "r", "https://cloudflare-eth.com", "specifies rpc url explicitly (shorthand)")
	flag.StringVar(&message, "message", "", "specifies message to be validated against")
	flag.StringVar(&message, "m", "", "specifies message to be validated against (shorthand)")
	flag.StringVar(&signer, "address", "", "specifies signer address")
	flag.StringVar(&signer, "a", "", "specifies signer address (shorthand)")
	flag.StringVar(&signature, "signature", "", "specifies signature to validate")
	flag.StringVar(&signature, "sig", "", "specifies signature to validate (shorthand)")
	flag.StringVar(&signature, "s", "", "specifies signature to validate (shorthand)")
	flag.StringVar(&validatorAddress, "validator", "", "specifies validator address (must be contract address)")
	flag.StringVar(&validatorAddress, "v", "", "specifies validator address (must be contract address) (shorthand)")
	flag.StringVar(&customValidSignature, "valid_signature", "", "specifies custom valid signature (successful response)")
	flag.StringVar(&customValidSignature, "vs", "", "specifies custom valid signature (successful response) (shorthand)")
	flag.BoolVar(&debug, "d", false, "enables debug comments (verbose)")

	flag.Parse()

	ctx := context.Background()
	logger := gaelogrus.GetLogger(ctx)

	if debug {
		logger.Logger.SetLevel(5)
	}

	if rpcURL == "" {
		logger.Error("empty rpc url provided")
		os.Exit(2)
	}

	if signer == "" {
		logger.Error("empty signer address provided")
		os.Exit(2)
	}

	if signature == "" {
		logger.Error("empty signature provided")
		os.Exit(2)
	}

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		logger.WithError(err).Fatalf("failed to validate signature")
	}

	logger.WithFields(map[string]interface{}{
		"signer":               signer,
		"message":              message,
		"signature":            signature,
		"rpcURL":               rpcURL,
		"validatorAddress":     validatorAddress,
		"customValidSignature": customValidSignature,
	}).Debug("arguments")

	validator := erc1271.NewValidator(client)

	if validatorAddress != "" {
		validator = validator.WithValidatorAddressHex(validatorAddress)
	}

	if customValidSignature != "" {
		validator = validator.WithCustomValidSignatureHex(customValidSignature)
	}

	valid, err := validator.Validate(
		ctx,
		[]byte(message),
		signer,
		signature,
	)
	if err != nil {
		logger.WithError(err).Fatalf("failed to validate signature")
	}

	logger.WithField("valid", valid).Info("result")
}
