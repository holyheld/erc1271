package erc1271

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, "https://cloudflare-eth.com")
	if err != nil {
		t.Fatal(err)
	}

	type Case struct {
		Description string
		Message     []byte
		Signer      string
		Signature   string
		Valid       bool
	}

	tests := []Case{
		{
			Description: "Valid ERC1271 signature",
			Message:     []byte("Hello go test!"),
			Signer:      "0x607377F587B1BDc68Bec3E19316D56bA8929d5eB",
			Signature:   "0xbcf08f9c64a93a58935c31e308b6e384cb72458e54cdb507a18c9b58fa7f910c6fb272c563e9aeaedf1e84fc8ecc8f2e840599ce612c8ade449004c1f575f89f1c",
			Valid:       true,
		},
		{
			Description: "Invalid ERC1271 signature",
			Message:     []byte("Hello go test!!"),
			Signer:      "0x607377F587B1BDc68Bec3E19316D56bA8929d5eB",
			Signature:   "0xbcf08f9c64a93a58935c31e308b6e384cb72458e54cdb507a18c9b58fa7f910c6fb272c563e9aeaedf1e84fc8ecc8f2e840599ce612c8ade449004c1f575f89f1c",
			Valid:       false,
		},
		{
			Description: "Valid EOA personal sign signature",
			Message:     []byte("Hello go test!"),
			Signer:      "0x5C6Aa53c883bB6c66CD2A0aD42Ae0828832A40E0",
			Signature:   "0x8a11e083dcb11229cec282c9b2122a05e1abbe9e4a6c98c04fb10537ac2585ab765d1bd31cd5c4bf194e55e27d7f2f498594fad9dc268787990337b6d46a71e51b",
			Valid:       false,
		},
		{
			Description: "Invalid EOA personal sign signature",
			Message:     []byte("Hello go test!!"),
			Signer:      "0x5C6Aa53c883bB6c66CD2A0aD42Ae0828832A40E0",
			Signature:   "0x8a11e083dcb11229cec282c9b2122a05e1abbe9e4a6c98c04fb10537ac2585ab765d1bd31cd5c4bf194e55e27d7f2f498594fad9dc268787990337b6d46a71e51b",
			Valid:       false,
		},
	}

	validator := NewERC1271Validator(client)
	for i, test := range tests {
		valid, err := validator.Validate(
			ctx,
			test.Message,
			test.Signer,
			test.Signature,
		)
		if err != nil {
			t.Errorf("%d (%s): expected err to be nil, got: %s", i, test.Description, err)
			continue
		}

		if !valid && test.Valid || valid && !test.Valid {
			t.Errorf("%d (%s): expected result to be %t, got: %t", i, test.Description, test.Valid, valid)
			continue
		}

		t.Logf("%d (%s): OK", i, test.Description)
	}
}
