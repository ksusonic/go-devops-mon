package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestServer_UpdateOneMetricHandler(t *testing.T) {
	var memStorage metrics.MetricStorage = storage.NewMemStorage()
	r := NewRouter(&memStorage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range []struct {
		MType          string
		MName          string
		MValue         float64
		ExpectedStatus int
	}{
		{
			"gauge",
			"BuckHashSys",
			123.01,
			http.StatusOK,
		},
		{
			"gauge",
			"noSuchMetric",
			124.01,
			http.StatusOK,
		},
		{
			"superGauge",
			"BuckHashSys",
			125.01,
			http.StatusNotImplemented,
		},
		{
			"superGauge",
			"BuckHashSys",
			125.01,
			http.StatusNotImplemented,
		},
	} {
		statusCode, _ := testRequest(t, ts, "POST", fmt.Sprintf("/update/%s/%s/%f", tt.MType, tt.MName, tt.MValue), nil)
		assert.Equal(t, tt.ExpectedStatus, statusCode)

		if tt.ExpectedStatus == http.StatusOK {
			// check how metric saved in storage
			actualMetric, err := memStorage.GetMetric(tt.MName)
			if err != nil {
				t.Errorf("metric %s not saved in storage", tt.MName)
			}
			assert.Equal(t, tt.MName, actualMetric.ID)
			assert.Equal(t, tt.MType, actualMetric.MType)
			assert.Equal(t, tt.MValue, *actualMetric.Value)
		}
	}

	statusCode, _ := testRequest(t, ts, "POST", "/update/counter/", nil)
	assert.Equal(t, http.StatusNotFound, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/noSuchCounter/1234567", nil)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestServer_UpdateMetricHandler(t *testing.T) {
	var memStorage metrics.MetricStorage = storage.NewMemStorage()
	r := NewRouter(&memStorage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var testValue = 123.0123

	for _, tt := range []struct {
		Metric         metrics.Metrics
		ExpectedStatus int
	}{
		{
			metrics.Metrics{
				ID:    "BuckHashSys",
				MType: "gauge",
				Delta: nil,
				Value: &testValue,
			},
			http.StatusOK,
		},
		{
			metrics.Metrics{
				ID:    "noSuchMetric",
				MType: "gauge",
				Delta: nil,
				Value: &testValue,
			},
			http.StatusOK,
		},
		{
			metrics.Metrics{
				ID:    "BuckHashSys",
				MType: "superGauge",
				Delta: nil,
				Value: &testValue,
			},
			http.StatusNotImplemented,
		},
	} {
		marshal, err := json.Marshal(tt.Metric)
		if err != nil {
			t.Error(err)
		}
		statusCode, _ := testRequest(t, ts, "POST", "/update/", bytes.NewReader(marshal))
		assert.Equal(t, tt.ExpectedStatus, statusCode)

		if tt.ExpectedStatus == http.StatusOK {
			// check how metric saved in storage
			actualMetric, err := memStorage.GetMetric(tt.Metric.ID)
			if err != nil {
				t.Errorf("metric %s not saved in storage", tt.Metric.ID)
			}
			assert.Equal(t, tt.Metric.ID, actualMetric.ID)
			assert.Equal(t, tt.Metric.MType, actualMetric.MType)
			assert.Equal(t, *tt.Metric.Value, *actualMetric.Value)
		}
	}

	var counterValue int64 = 1234
	var simpleCounterMetric = metrics.Metrics{
		ID:    "noSuchCounter",
		MType: "counter",
		Delta: &counterValue,
		Value: nil,
	}
	marshal, _ := json.Marshal(simpleCounterMetric)

	statusCode, _ := testRequest(t, ts, "POST", "/update/counter/", bytes.NewReader(marshal))
	assert.Equal(t, http.StatusNotFound, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/noSuchCounter/1234567", bytes.NewReader(marshal))
	assert.Equal(t, http.StatusOK, statusCode)
}
