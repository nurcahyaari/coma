package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

func readRSAPublicKey() *rsa.PublicKey {
	file, err := os.ReadFile(CONST.DEFAULT_RSA_PUBLIC_KEY_LOCATION)
	if err != nil {
		log.Fatal().Err(err)
	}

	block, _ := pem.Decode(file)
	if block == nil {
		log.Fatal().Err(errors.New("err: cannot decode file"))
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Fatal().Err(err)
	}

	return publicKey
}

func readRSAPrivateKey() *rsa.PrivateKey {
	file, err := os.ReadFile(CONST.DEFAULT_RSA_PRIVATE_KEY_LOCATION)
	if err != nil {
		log.Fatal().Err(err)
	}

	block, _ := pem.Decode(file)
	if block == nil {
		log.Fatal().Err(errors.New("err: cannot decode file"))
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal().Err(err)
	}

	return privateKey
}

func createDefaultRSAKey() {

}
