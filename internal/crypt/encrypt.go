package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Encrypter struct {
	publicKey *rsa.PublicKey
}

func NewEncrypter(publicKeyPath string) (*Encrypter, error) {
	bytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	keyPem, _ := pem.Decode(bytes)
	if keyPem == nil {
		return nil, errors.New("incorrect public key")
	}
	cert, err := x509.ParseCertificate(keyPem.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := cert.PublicKey.(*rsa.PublicKey)

	return &Encrypter{publicKey: publicKey}, nil
}

func (e *Encrypter) EncryptBytes(b []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, b)
}
