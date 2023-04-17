package erc1271

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()

	clients := make(map[int]*ethclient.Client, 2)

	client, err := ethclient.DialContext(ctx, "https://cloudflare-eth.com")
	if err != nil {
		t.Fatal(err)
	}
	clients[1] = client

	client, err = ethclient.DialContext(ctx, "https://polygon-rpc.com")
	if err != nil {
		t.Fatal(err)
	}
	clients[137] = client

	type Case struct {
		Description string
		Message     []byte
		Signer      string
		Signature   string
		Valid       bool
		ChainID     int
	}

	tests := []Case{
		{
			Description: "Valid ERC1271 signature (Argent wallet)",
			Message:     []byte("Hello go test!"),
			Signer:      "0x607377F587B1BDc68Bec3E19316D56bA8929d5eB",
			Signature:   "0xbcf08f9c64a93a58935c31e308b6e384cb72458e54cdb507a18c9b58fa7f910c6fb272c563e9aeaedf1e84fc8ecc8f2e840599ce612c8ade449004c1f575f89f1c",
			Valid:       true,
			ChainID:     1,
		},
		{
			Description: "Valid ERC1271 signature (Ambire wallet)",
			Message:     []byte("Hello go test!"),
			Signer:      "0xAB833C2DDb1394Cf14AAdCc0aCa4B66Ee84d2C74",
			Signature:   "0x000000000000000000000000000000000000000000000000000000000003f480000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000042c44821b4f6bccb7e1ebe530766f6713df0c12a02b135e4bdc36b7ce3404603e56f511ef5341e2cd123c4101debf8d8271677caa0dda68fc8875bcf83282256c31b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004281ca98f784a73d5e623de6c3992bf78f8b3565bd4cc29df7dbe52473334bab464e76c597bcf7caebae9bdb3b1ef9708160210bf78863b747c39d0b0763f763441c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000ff3f6d14df43c112ab98834ee1f82083e07c26bf02",
			Valid:       true,
			ChainID:     137,
		},
		{
			Description: "Invalid ERC1271 signature (Ambire wallet)",
			Message:     []byte("Hello go test!!"),
			Signer:      "0xAB833C2DDb1394Cf14AAdCc0aCa4B66Ee84d2C74",
			Signature:   "0x000000000000000000000000000000000000000000000000000000000003f480000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000042c44821b4f6bccb7e1ebe530766f6713df0c12a02b135e4bdc36b7ce3404603e56f511ef5341e2cd123c4101debf8d8271677caa0dda68fc8875bcf83282256c31b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004281ca98f784a73d5e623de6c3992bf78f8b3565bd4cc29df7dbe52473334bab464e76c597bcf7caebae9bdb3b1ef9708160210bf78863b747c39d0b0763f763441c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000ff3f6d14df43c112ab98834ee1f82083e07c26bf02",
			Valid:       false,
			ChainID:     137,
		},
		{
			Description: "Invalid ERC1271 signature",
			Message:     []byte("Hello go test!!"),
			Signer:      "0x607377F587B1BDc68Bec3E19316D56bA8929d5eB",
			Signature:   "0xbcf08f9c64a93a58935c31e308b6e384cb72458e54cdb507a18c9b58fa7f910c6fb272c563e9aeaedf1e84fc8ecc8f2e840599ce612c8ade449004c1f575f89f1c",
			Valid:       false,
			ChainID:     1,
		},
		{
			Description: "Valid EOA personal sign signature",
			Message:     []byte("Hello go test!"),
			Signer:      "0x5C6Aa53c883bB6c66CD2A0aD42Ae0828832A40E0",
			Signature:   "0x8a11e083dcb11229cec282c9b2122a05e1abbe9e4a6c98c04fb10537ac2585ab765d1bd31cd5c4bf194e55e27d7f2f498594fad9dc268787990337b6d46a71e51b",
			Valid:       false,
			ChainID:     1,
		},
		{
			Description: "Invalid EOA personal sign signature",
			Message:     []byte("Hello go test!!"),
			Signer:      "0x5C6Aa53c883bB6c66CD2A0aD42Ae0828832A40E0",
			Signature:   "0x8a11e083dcb11229cec282c9b2122a05e1abbe9e4a6c98c04fb10537ac2585ab765d1bd31cd5c4bf194e55e27d7f2f498594fad9dc268787990337b6d46a71e51b",
			Valid:       false,
			ChainID:     1,
		},
	}

	for i, test := range tests {
		client, ok := clients[test.ChainID]
		if !ok {
			t.Errorf("%d (%s): could not find client with chain id of %d", i, test.Description, test.ChainID)
			continue
		}
		valid, err := NewValidator(client).Validate(
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
