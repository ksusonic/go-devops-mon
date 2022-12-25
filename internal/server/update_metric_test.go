package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestServer(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/BuckHashSys/123.01")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/gauge/noSuchMetric/123.01")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/superGauge/BuckHashSys/123.01")
	assert.Equal(t, http.StatusNotImplemented, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/")
	assert.Equal(t, http.StatusNotFound, statusCode)
}
