package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSA(t *testing.T) {

	GenerateKeyPair := func() (*rsa.PrivateKey, *rsa.PublicKey) {
		privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err, "Cannot generate RSA key")
		return privatekey, &privatekey.PublicKey
	}

	privateKey, publicKey := GenerateKeyPair()

	t.Run("dump private key to private.pem", func(t *testing.T) {

		var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		}

		privatePemFile, err := os.Create("testdata/private.pem")
		require.NoError(t, err, "Cannot create private.pem")

		err = pem.Encode(privatePemFile, privateKeyBlock)
		require.NoError(t, err, "Cannot write to private.pem")

		GetPrivateKey := func() *rsa.PrivateKey {
			b, err := os.ReadFile("testdata/private.pem")
			assert.NoError(t, err)
	
			block, _ := pem.Decode(b)
			assert.NotZero(t, block)
	
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			assert.NoError(t, err)
	
			return key
		}

		assert.True(t, privateKey.Equal(GetPrivateKey()))
	})

	t.Run("dump public key to public.pem", func(t *testing.T) {

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		require.NoError(t, err, "Cannot dumping publickey")

		publicKeyBlock := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}

		publicPem, err := os.Create("testdata/public.pem")
		require.NoError(t, err, "Cannot create public.pem")

		err = pem.Encode(publicPem, publicKeyBlock)
		require.NoError(t, err, "Cannot write to public.pem")

		GetPublicKey := func() *rsa.PublicKey {
			b, err := os.ReadFile("testdata/public.pem")
			assert.NoError(t, err)

			block, _ := pem.Decode(b)
			assert.NotZero(t, block)

			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			assert.NoError(t, err)
			return key.(*rsa.PublicKey)
		}

		assert.True(t, publicKey.Equal(GetPublicKey()))
	})

}
