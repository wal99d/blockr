package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {

	privKey := GeneratePrivateKey()

	assert.Equal(t, len(privKey.Bytes()), privKeyLen)
	pubkey := privKey.Public()

	assert.Equal(t, len(pubkey.Bytes()), pubKeyLen)
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		addressStr = "701518c6d2aa32e3476007d1e36670f5518a6fe8"
		seed       = "18eb974804b28c982fd67c71d4c78cd76bd6f0505ed7f5e7c1a8efb7df2744bb"
		privKey    = NewPrivateKeyFromString(seed)
	)
	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	address := privKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))
	//test with invlaid msg
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	//test with invalid pubkey
	invalidPrivKey := GeneratePrivateKey()
	invalidPubKey := invalidPrivKey.Public()
	assert.False(t, sig.Verify(invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {

	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, addressLen, len(address.Bytes()))
	fmt.Println(address)
}
