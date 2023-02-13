package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type Service struct {
	key *string
}

func NewService(key string) *Service {
	var keyPtr *string
	if key != "" {
		keyPtr = &key
	}
	return &Service{key: keyPtr}
}

func (s Service) SetHash(m *metrics.Metrics) error {
	if s.key == nil {
		return nil
	}

	hash, err := s.hash(*s.key, m)
	if err != nil {
		return err
	}
	m.Hash = hash
	return nil
}

func (s Service) ValidateHash(m *metrics.Metrics) error {
	if s.key == nil {
		return nil
	}

	hash, err := s.hash(*s.key, m)
	if err != nil {
		return err
	}

	if hash != m.Hash {
		return fmt.Errorf("hash is not correct:\nexpected:'%s'\nactual:'%s'", hash, m.Hash)
	}
	return nil
}

func (s Service) hash(key string, m *metrics.Metrics) (hash string, err error) {
	switch m.MType {
	case metrics.CounterMType:
		hash = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case metrics.GaugeMType:
		hash = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		return "", fmt.Errorf("cannot calc hash of type %s", m.MType)
	}
	h := hmac.New(sha256.New, []byte(key))
	_, err = h.Write([]byte(hash))
	if err != nil {
		return "", fmt.Errorf("cannot calc hash: %v", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
