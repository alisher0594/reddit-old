package main

import (
	"github.com/alisher0594/reddit-old/internal/data/models"
	"github.com/alisher0594/reddit-old/internal/jsonlog"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *application {
	return &application{
		config: config{
			port: 4000,
			env:  "staging",
		},
		logger: jsonlog.New(ioutil.Discard, jsonlog.LevelOff),
		models: models.NewMock(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	return &testServer{httptest.NewTLSServer(h)}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
