package main_test

import (
	"bytes"
	"io"
	"kvstore/cmd/web"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("Wanted a status code of %d but got %d", http.StatusOK, code)
	}

	if body != "OK" {
		t.Errorf("Wanted a body of \"OK\" but got %q", body)
	}

}

func TestNoSuchEndpointReturnsNotFound(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(app.Routes())
	defer ts.Close()

	got, _, _ := ts.get(t, "/nosuchroute")
	want := http.StatusNotFound

	if got != want {
		t.Errorf("Wanted a status code of %d but got %d", want, got)
	}

}

func newTestApplication() *main.Application {

	return &main.Application{
		ErrorLog: log.New(io.Discard, "", 0),
		InfoLog:  log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	t.Helper()

	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
