package crypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type Decrypter struct {
	privateKey *rsa.PrivateKey
	logger     *zap.Logger
}

func NewDecrypter(privateKeyPath string, logger *zap.Logger) (*Decrypter, error) {
	file, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKeyPem, _ := pem.Decode(file)
	if privateKeyPem == nil {
		return nil, errors.New("incorrect private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
	if err != nil {
		return nil, err
	}

	return &Decrypter{
		privateKey: key,
		logger:     logger,
	}, nil
}

func (d *Decrypter) DecryptBytes(b []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, d.privateKey, b)
}

func (d *Decrypter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		all, err := io.ReadAll(r.Body)
		if err != nil {
			d.logger.Error("error reading request body", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		decryptedBytes, err := d.DecryptBytes(all)
		if err != nil {
			d.logger.Error("error reading request body", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(decryptedBytes))
		next.ServeHTTP(rw, r)
	})
}
