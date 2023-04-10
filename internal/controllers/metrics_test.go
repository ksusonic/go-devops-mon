package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/hash"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type TestServer struct {
	Server *httptest.Server
}

func NewTestServer(metricStorage metrics.ServerMetricStorage) TestServer {
	router := chi.NewRouter()
	r := NewMetricController(zap.NewExample(), metricStorage, hash.NewService(""))
	router.Mount("/", r.Router())
	ts := httptest.NewServer(router)
	return TestServer{
		Server: ts,
	}
}

func (s *TestServer) Close() {
	s.Server.Close()
}

func (s *TestServer) testRequest(t *testing.T, method, path string, body io.Reader) (int, []byte) {
	req, err := http.NewRequest(method, s.Server.URL+path, body)
	require.NoError(t, err)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, respBody
}

func TestController_UpdateMetricPathHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	ts := NewTestServer(memStorage)
	defer ts.Close()

	statusCode, _ := ts.testRequest(t, "POST", "/update/gauge/BuckHashSys/123.01", nil)
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = ts.testRequest(t, "POST", "/update/gauge/noSuchMetric/123.01", nil)
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = ts.testRequest(t, "POST", "/update/superGauge/BuckHashSys/123.01", nil)
	assert.Equal(t, http.StatusNotImplemented, statusCode)

	statusCode, _ = ts.testRequest(t, "POST", "/update/counter/", nil)
	assert.Equal(t, http.StatusNotFound, statusCode)

	statusCode, _ = ts.testRequest(t, "POST", "/update/counter/RandomValue/12345678", nil)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestController_UpdateMetricHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	ts := NewTestServer(memStorage)
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
		statusCode, _ := ts.testRequest(t, "POST", "/update/", bytes.NewReader(marshal))
		assert.Equal(t, tt.ExpectedStatus, statusCode)

		if tt.ExpectedStatus == http.StatusOK {
			// check how metric saved in storage
			actualMetric, err := memStorage.GetMetric(context.Background(), tt.Metric.MType, tt.Metric.ID)
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

	statusCode, _ := ts.testRequest(t, "POST", "/update/counter/", bytes.NewReader(marshal))
	assert.Equal(t, http.StatusNotFound, statusCode)

	statusCode, _ = ts.testRequest(t, "POST", "/update/counter/noSuchCounter/1234567", bytes.NewReader(marshal))
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestController_GetMetricPathHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	ts := NewTestServer(memStorage)
	defer ts.Close()

	statusCode, _ := ts.testRequest(t, "POST", "/update/gauge/BuckHashSys/123.01", nil)
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, res := ts.testRequest(t, "GET", "/value/gauge/BuckHashSys", nil)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "123.01", string(res))
}

func TestController_GetMetricHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	ts := NewTestServer(memStorage)
	defer ts.Close()

	testValue := 123.01
	expected := metrics.Metrics{
		ID:    "BuckHashSys",
		MType: "gauge",
		Delta: nil,
		Value: &testValue,
	}
	marshal, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}
	statusCode, _ := ts.testRequest(t, "POST", "/update/", bytes.NewReader(marshal))
	assert.Equal(t, http.StatusOK, statusCode)

	marshal, err = json.Marshal(metrics.Metrics{
		ID:    "BuckHashSys",
		MType: "gauge",
	})
	if err != nil {
		t.Error(err)
	}
	statusCode, res := ts.testRequest(t, "POST", "/value/", bytes.NewReader(marshal))
	assert.Equal(t, http.StatusOK, statusCode)
	var actual metrics.Metrics
	err = json.Unmarshal(res, &actual)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, actual)
}

func ExampleController_GetMetricHandler() {
	memStorage := storage.NewMemStorage()
	testValue := 123.01
	expected := metrics.Metrics{
		ID:    "BuckHashSys",
		MType: "gauge",
		Delta: nil,
		Value: &testValue,
	}
	memStorage.SetMetric(context.Background(), expected)
	ts := NewTestServer(memStorage)
	defer ts.Close()

	marshal, _ := json.Marshal(metrics.Metrics{
		ID:    "BuckHashSys",
		MType: "gauge",
	})
	req, _ := http.NewRequest("POST", ts.Server.URL+"/value/", bytes.NewReader(marshal))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var actual metrics.Metrics
	json.Unmarshal(respBody, &actual)

	fmt.Println(actual)
	actual.Hash = "some-hash"
	fmt.Println(actual)

	// Output:
	// metric BuckHashSys of type gauge with value 123.010000
	// metric BuckHashSys of type gauge with value 123.010000 and hash: some-hash
}

func TestController_GetAllMetricsHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	ts := NewTestServer(memStorage)
	defer ts.Close()

	for _, request := range []string{
		"/update/gauge/HeapInuse/786432.01",
		"/update/gauge/HeapObjects/613",
		"/update/gauge/NextGC/4194304",
		"/update/gauge/LastGC/0",
		"/update/gauge/MSpanInuse/30528",
		"/update/gauge/StackInuse/327680",
		"/update/gauge/GCCPUFraction/0",
		"/update/gauge/Lookups/0",
		"/update/gauge/MCacheSys/15600",
		"/update/gauge/PauseTotalNs/0",
		"/update/gauge/NumGC/0",
		"/update/gauge/OtherSys/749387",
		"/update/gauge/GCSys/7963072",
		"/update/gauge/HeapIdle/3080192",
		"/update/gauge/HeapReleased/3047424",
		"/update/gauge/MCacheInuse/14400",
		"/update/gauge/Alloc/218216",
		"/update/gauge/MSpanSys/32544",
		"/update/gauge/Sys/12958736",
		"/update/gauge/BuckHashSys/3829",
		"/update/gauge/Frees/24",
		"/update/gauge/Mallocs/637",
		"/update/gauge/NumForcedGC/0",
		"/update/gauge/HeapAlloc/218216.123",
		"/update/gauge/HeapSys/3866624",
		"/update/gauge/StackSys/327680",
		"/update/gauge/TotalAlloc/218216",
		"/update/counter/PollCount/5",
		"/update/counter/RandomValue/3916589616287113937",
	} {
		statusCode, _ := ts.testRequest(t, "POST", request, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	}
	expected := `{"counter":{"PollCount":5,"RandomValue":3916589616287113937},"gauge":{"Alloc":218216,"BuckHashSys":3829,"Frees":24,"GCCPUFraction":0,"GCSys":7963072,"HeapAlloc":218216.123,"HeapIdle":3080192,"HeapInuse":786432.01,"HeapObjects":613,"HeapReleased":3047424,"HeapSys":3866624,"LastGC":0,"Lookups":0,"MCacheInuse":14400,"MCacheSys":15600,"MSpanInuse":30528,"MSpanSys":32544,"Mallocs":637,"NextGC":4194304,"NumForcedGC":0,"NumGC":0,"OtherSys":749387,"PauseTotalNs":0,"StackInuse":327680,"StackSys":327680,"Sys":12958736,"TotalAlloc":218216}}`

	statusCode, actual := ts.testRequest(t, "GET", "/", nil)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.JSONEq(t, expected, string(actual))
}
