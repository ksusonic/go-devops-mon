package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_GetMetric(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/BuckHashSys/123.01")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, res := testRequest(t, ts, "GET", "/value/gauge/BuckHashSys")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "123.01", res)
}
