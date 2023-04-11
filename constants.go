package erc1271

import "github.com/ethereum/go-ethereum/crypto"

// ValidSignature is a magic value to compare validate result against
var ValidSignature = crypto.Keccak256([]byte("isValidSignature(bytes32,bytes)"))[:4]
