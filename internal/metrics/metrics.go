package metrics

import (
	"crypto/sha256"
	"errors"
	"fmt"
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
		return errors.New("missing ID of metric")
	} else if m.MType == "" {
		return errors.New("missing type of metric")
	}

	return nil
}

func (m Metrics) CalcHash(key string) string {
	var hash string
	fmt.Println("metric - ", m.ID, m.Delta, m.Value)
	switch m.MType {
	case CounterMType:
		if m.Delta == nil {
			return ""
		}
		hash = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case GaugeMType:
		hash = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		return ""
	}
	h := sha256.New()
	h.Write([]byte(hash))
	h.Write([]byte(key))
	return string(h.Sum(nil))
}

func (m Metrics) ValidateHash(key string) error {
	if m.Hash == "" {
		return nil
	}

	if m.CalcHash(key) != m.Hash {
		return errors.New("hash is not correct")
	}
	return nil
}

func (m Metrics) String() string {
	switch m.MType {
	case CounterMType:
		return fmt.Sprintf("metric %s of type %s with value %d", m.ID, m.MType, *m.Delta)
	case GaugeMType:
		return fmt.Sprintf("metric %s of type %s with value %f", m.ID, m.MType, *m.Value)
	default:
		return ""
	}
}
