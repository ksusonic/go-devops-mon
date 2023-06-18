package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	metricspb "github.com/ksusonic/go-devops-mon/proto/metrics"
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

func (s Service) SetHashProto(m *metricspb.Metric) error {
	if s.key == nil {
		return nil
	}

	h, err := hash(*s.key, m)
	if err != nil {
		return err
	}
	m.Hash = h
	return nil
}

func (s Service) SetHash(m *metrics.Metric) error {
	if s.key == nil {
		return nil
	}

	h, err := hash(*s.key, m)
	if err != nil {
		return err
	}
	m.Hash = h
	return nil
}

func (s Service) ValidateHash(m *metrics.Metric) error {
	if s.key == nil {
		return nil
	}

	hash, err := hash(*s.key, m)
	if err != nil {
		return err
	}

	if hash != m.Hash {
		return fmt.Errorf("hash is not correct:\nexpected:'%s'\nactual:'%s'", hash, m.Hash)
	}
	return nil
}

func (s Service) ValidateHashProto(m *metricspb.Metric) error {
	if s.key == nil {
		return nil
	}

	hash, err := hash(*s.key, m)
	if err != nil {
		return err
	}

	if hash != m.Hash {
		return fmt.Errorf("hash is not correct:\nexpected:'%s'\nactual:'%s'", hash, m.Hash)
	}
	return nil
}

func hash[M metrics.GenericMetric](key string, m M) (hash string, err error) {
	switch m.GetType() {
	case metricspb.MetricType_counter:
		hash = fmt.Sprintf("%s:counter:%d", m.GetID(), m.GetDelta())
	case metricspb.MetricType_gauge:
		hash = fmt.Sprintf("%s:gauge:%f", m.GetID(), m.GetValue())
	default:
		return "", fmt.Errorf("cannot calc hash of type %s", m.GetID())
	}
	h := hmac.New(sha256.New, []byte(key))
	_, err = h.Write([]byte(hash))
	if err != nil {
		return "", fmt.Errorf("cannot calc hash: %v", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
