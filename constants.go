package erc1271

import "github.com/ethereum/go-ethereum/crypto"

var validSignature = crypto.Keccak256([]byte("isValidSignature(bytes32,bytes)"))[:4]
