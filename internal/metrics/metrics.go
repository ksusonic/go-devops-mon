package metrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
)

const (
	GaugeMType   = "gauge"
	CounterMType = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (m Metrics) Bind(*http.Request) error {
	if m.ID == "" {
		return fmt.Errorf("missing ID of metric")
	} else if m.MType == "" {
		return fmt.Errorf("missing type of metric")
	}

	return nil
}

func (m Metrics) CalcHash(key string) (string, error) {
	var hash string
	switch m.MType {
	case CounterMType:
		hash = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case GaugeMType:
		hash = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		return "", fmt.Errorf("cannot calc hash of type %s", m.MType)
	}
	log.Printf("used %s + %s for hash\n", hash, key)
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(hash))
	if err != nil {
		return "", fmt.Errorf("cannot calc hash: %v", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (m Metrics) ValidateHash(key string) error {
	calculated, err := m.CalcHash(key)
	if err != nil {
		return fmt.Errorf("cannot calc hash: %v", err)
	}

	if calculated != m.Hash {
		return fmt.Errorf("hash does not match: actual: %s expected: %s", calculated, m.Hash)
	}
	return nil
}

func (m Metrics) String() string {
	switch m.MType {
	case CounterMType:
		return fmt.Sprintf("metric %s of type %s with value %d and hash %s", m.ID, m.MType, *m.Delta, m.Hash)
	case GaugeMType:
		return fmt.Sprintf("metric %s of type %s with value %f and hash %s", m.ID, m.MType, *m.Value, m.Hash)
	default:
		return ""
	}
}
