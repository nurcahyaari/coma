package config

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"golang.org/x/sync/errgroup"
)

const (
	PRIVATE_KEY_CODE = "RSA PRIVATE KEY"
	PUBLIC_KEY_CODE  = "PUBLIC KEY"
)

func createDefaultRSAIfNotExist() error {
	if _, err := os.Stat(CONST.DEFAULT_RSA_PRIVATE_KEY_LOCATION); !os.IsNotExist(err) {
		return nil
	}

	privateRsaKey, err := CreateDefaultRSAPrivateKey()
	if err != nil {
		return err
	}

	if privateRsaKey == nil {
		return errors.New("error: cannot create rsa")
	}

	privateRsaKeyBytes := encodePrivateKeyToPEM(privateRsaKey)
	publicRsaKeyBytes, err := encodePublicKeyToPEM(privateRsaKey)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.TODO())

	// save private key
	group.Go(func() error {
		return os.WriteFile(CONST.DEFAULT_RSA_PRIVATE_KEY_LOCATION, privateRsaKeyBytes, 0600)
	})

	// save public key
	group.Go(func() error {
		return os.WriteFile(CONST.DEFAULT_RSA_PUBLIC_KEY_LOCATION, publicRsaKeyBytes, 0600)
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func CreateDefaultRSAPrivateKey() (*rsa.PrivateKey, error) {
	bitSize := 1028
	if CONST.DEFAULT_RSA_BITSIZE != 0 {
		bitSize = CONST.DEFAULT_RSA_BITSIZE
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  PRIVATE_KEY_CODE,
		Bytes: privDER,
	})

	return privatePEM
}

func encodePublicKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  PUBLIC_KEY_CODE,
		Bytes: pubDER,
	})

	return publicPEM, nil
}

func readRSAPublicKey() *rsa.PublicKey {
	file, err := os.ReadFile(CONST.DEFAULT_RSA_PUBLIC_KEY_LOCATION)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(file)

	if block == nil || block.Type != PUBLIC_KEY_CODE {
		panic(errors.New("err: cannot decode file"))
	}

	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	publicKey, ok := public.(*rsa.PublicKey)
	if !ok {
		panic(errors.New("not an RSA public key"))
	}

	return publicKey
}

func readRSAPrivateKey() *rsa.PrivateKey {
	file, err := os.ReadFile(CONST.DEFAULT_RSA_PRIVATE_KEY_LOCATION)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(file)
	if block == nil || block.Type != PRIVATE_KEY_CODE {
		panic(errors.New("err: cannot decode file"))
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	return privateKey
}
