package server

import (
	"github.com/ksusonic/go-devops-mon/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_UpdateMetric(t *testing.T) {
	s := Server{MemStorage: &storage.MemStorage{
		GaugeStorage:   storage.GaugeStorage{},
		CounterStorage: storage.CounterStorage{},
	}}

	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "simple test #1",
			request: "/update/gauge/BuckHashSys/123.01",
			want: want{
				statusCode: 200,
			},
		},
		{
			name:    "unknown metric name test #2",
			request: "/update/gauge/noSuchMetric/123.01",
			want: want{
				statusCode: 404,
			},
		},
		{
			name:    "unknown metric type test #2",
			request: "/update/superGauge/BuckHashSys/123.01",
			want: want{
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(s.UpdateMetric)
			h(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
