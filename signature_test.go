package crypto

import (
	"strings"
	"testing"

	"github.com/neatio-net/ed25519"
	data "github.com/neatio-net/data-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignAndValidateEd25519(t *testing.T) {

	privKey := GenPrivKeyEd25519()
	pubKey := privKey.PubKey()

	msg := CRandBytes(128)
	sig := privKey.Sign(msg)

	assert.True(t, pubKey.VerifyBytes(msg, sig))

	sigEd := sig.(SignatureEd25519)
	sigEd[0] ^= byte(0x01)
	sig = Signature(sigEd)

	assert.False(t, pubKey.VerifyBytes(msg, sig))
}

func TestSignAndValidateSecp256k1(t *testing.T) {
	privKey := GenPrivKeySecp256k1()
	pubKey := privKey.PubKey()

	msg := CRandBytes(128)
	sig := privKey.Sign(msg)

	assert.True(t, pubKey.VerifyBytes(msg, sig))

	sigEd := sig.(SignatureSecp256k1)
	sigEd[0] ^= byte(0x01)
	sig = Signature(sigEd)

	assert.False(t, pubKey.VerifyBytes(msg, sig))
}

func TestSignatureEncodings(t *testing.T) {
	cases := []struct {
		privKey PrivKeyS
		sigSize int
		sigType byte
		sigName string
	}{
		{
			privKey: PrivKeyS{GenPrivKeyEd25519()},
			sigSize: ed25519.SignatureSize,
			sigType: TypeEd25519,
			sigName: NameEd25519,
		},
		{
			privKey: PrivKeyS{GenPrivKeySecp256k1()},
			sigSize: 0,
			sigType: TypeSecp256k1,
			sigName: NameSecp256k1,
		},
	}

	for _, tc := range cases {

		pubKey := PubKeyS{tc.privKey.PubKey()}

		msg := CRandBytes(128)
		sig := SignatureS{tc.privKey.Sign(msg)}

		bin, err := data.ToWire(sig)
		require.Nil(t, err, "%+v", err)
		if tc.sigSize != 0 {
			assert.Equal(t, tc.sigSize+1, len(bin))
		}
		assert.Equal(t, tc.sigType, bin[0])

		sig2 := SignatureS{}
		err = data.FromWire(bin, &sig2)
		require.Nil(t, err, "%+v", err)
		assert.EqualValues(t, sig, sig2)
		assert.True(t, pubKey.VerifyBytes(msg, sig2))

		js, err := data.ToJSON(sig)
		require.Nil(t, err, "%+v", err)
		assert.True(t, strings.Contains(string(js), tc.sigName))

		sig3 := SignatureS{}
		err = data.FromJSON(js, &sig3)
		require.Nil(t, err, "%+v", err)
		assert.EqualValues(t, sig, sig3)
		assert.True(t, pubKey.VerifyBytes(msg, sig3))

		text, err := data.ToText(sig)
		require.Nil(t, err, "%+v", err)
		assert.True(t, strings.HasPrefix(text, tc.sigName))
	}
}

func TestWrapping(t *testing.T) {
	assert := assert.New(t)

	msg := CRandBytes(128)
	priv := GenPrivKeyEd25519()
	pub := priv.PubKey()
	sig := priv.Sign(msg)

	pubs := []PubKeyS{
		WrapPubKey(nil),
		WrapPubKey(pub),
		WrapPubKey(WrapPubKey(WrapPubKey(WrapPubKey(pub)))),
		WrapPubKey(PubKeyS{PubKeyS{PubKeyS{pub}}}),
	}
	for _, p := range pubs {
		_, ok := p.PubKey.(PubKeyS)
		assert.False(ok)
	}

	sigs := []SignatureS{
		WrapSignature(nil),
		WrapSignature(sig),
		WrapSignature(WrapSignature(WrapSignature(WrapSignature(sig)))),
		WrapSignature(SignatureS{SignatureS{SignatureS{sig}}}),
	}
	for _, s := range sigs {
		_, ok := s.Signature.(SignatureS)
		assert.False(ok)
	}

}
