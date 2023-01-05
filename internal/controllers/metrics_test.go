package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestServer_UpdateMetric(t *testing.T) {
	var memStorage metrics.ServerMetricStorage = storage.NewMemStorage()
	router := chi.NewRouter()
	r := NewMetricController(memStorage)
	router.Mount("/", r.Router())
	ts := httptest.NewServer(router)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/BuckHashSys/123.01")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/gauge/noSuchMetric/123.01")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/superGauge/BuckHashSys/123.01")
	assert.Equal(t, http.StatusNotImplemented, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/")
	assert.Equal(t, http.StatusNotFound, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/RandomValue/12345678")
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestServer_GetMetric(t *testing.T) {
	var memStorage metrics.ServerMetricStorage = storage.NewMemStorage()
	router := chi.NewRouter()
	r := NewMetricController(memStorage)
	router.Mount("/", r.Router())
	ts := httptest.NewServer(router)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/BuckHashSys/123.01", nil)
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, res := testRequest(t, ts, "GET", "/value/gauge/BuckHashSys", nil)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "123.01", res)
}

func TestServer_GetAllMetrics(t *testing.T) {
	var memStorage metrics.ServerMetricStorage = storage.NewMemStorage()
	router := chi.NewRouter()
	r := NewMetricController(memStorage)
	router.Mount("/", r.Router())
	ts := httptest.NewServer(router)
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
		statusCode, _ := testRequest(t, ts, "POST", request, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	}
	expected := `{"counter":{"PollCount":5,"RandomValue":3916589616287113937},"gauge":{"Alloc":218216,"BuckHashSys":3829,"Frees":24,"GCCPUFraction":0,"GCSys":7963072,"HeapAlloc":218216.123,"HeapIdle":3080192,"HeapInuse":786432.01,"HeapObjects":613,"HeapReleased":3047424,"HeapSys":3866624,"LastGC":0,"Lookups":0,"MCacheInuse":14400,"MCacheSys":15600,"MSpanInuse":30528,"MSpanSys":32544,"Mallocs":637,"NextGC":4194304,"NumForcedGC":0,"NumGC":0,"OtherSys":749387,"PauseTotalNs":0,"StackInuse":327680,"StackSys":327680,"Sys":12958736,"TotalAlloc":218216}}`

	statusCode, actual := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.JSONEq(t, expected, actual)
}
