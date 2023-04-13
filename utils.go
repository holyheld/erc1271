package erc1271

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

// IsZeroAddress validates if it's a 0 address
//
// https://github.com/miguelmota/ethereum-development-with-go-book/blob/3b24cad5d54aff1496fe4c3590c27f0fad648ddf/code/util/util.go#L42
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}
