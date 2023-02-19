package main_test

import (
	"bytes"
	"encoding/json"
	"io"
	"kvstore/cmd/web"
	"kvstore/kv"
	"kvstore/kv/memory"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	t.Parallel()
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
	t.Parallel()
	app := newTestApplication()

	ts := newTestServer(app.Routes())
	defer ts.Close()

	got, _, _ := ts.get(t, "/nosuchroute")
	want := http.StatusNotFound

	if got != want {
		t.Errorf("Wanted a status code of %d but got %d", want, got)
	}

}

func TestStoreValueWithValidInput(t *testing.T) {

	t.Parallel()
	got := make(map[string]string)
	store := &simpleStore{store: got}
	app := newTestApplication(withStore(store))

	ts := newTestServer(app.Routes())
	defer ts.Close()

	wantKey := "foo"
	wantValue := "bar"

	rc, _, _ := ts.put(t, "/store/"+wantKey, valueToJSON(t, wantValue))

	if rc != http.StatusOK {
		t.Errorf("Wanted a status code of %d but got %d", http.StatusOK, rc)
	}

	value, found := got[wantKey]

	if !found {
		t.Fatalf("No value set against our key %q", wantKey)
	}

	if value != wantValue {
		t.Fatalf("Value set against our key should be %q but is %q", wantKey, wantValue)
	}

}

func TestStoreValueWithEmptyKey(t *testing.T) {

	t.Parallel()

	app := newTestApplication()

	ts := newTestServer(app.Routes())
	defer ts.Close()

	wantKey := ""
	wantValue := "bar"

	rc, _, _ := ts.put(t, "/set/"+wantKey, valueToJSON(t, wantValue))

	// TODO: should maybe be a bad request instead ?
	want := http.StatusNotFound

	if rc != want {
		t.Errorf("Wanted a status code of %d but got %d", want, rc)
	}

}

func TestGetValueSetAgainstExistingKey(t *testing.T) {

	t.Parallel()

	wantKey := "previously_set"
	wantValue := "expected_value"
	store := &simpleStore{store: map[string]string{wantKey: wantValue}}
	app := newTestApplication(withStore(store))

	ts := newTestServer(app.Routes())
	defer ts.Close()

	rc, header, data := ts.get(t, "/store/"+wantKey)

	if rc != http.StatusOK {
		t.Errorf("Wanted a status code of %d but got %d", http.StatusOK, rc)
	}

	encoding := header.Get("Content-Type")

	if encoding != "application/json" {
		t.Fatalf("Expecting content type of application/json, but got %q", encoding)
	}

	got := valueFromJSON(t, data)

	if got != wantValue {
		t.Fatalf("Expecting a value of %q but got %q", wantValue, got)
	}

}

func newTestApplication(opts ...func(application *main.Application)) *main.Application {

	app := &main.Application{
		ErrorLog: log.New(io.Discard, "", 0),
		InfoLog:  log.New(io.Discard, "", 0),
		Store:    memory.NewStore(),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app

}

func withStore(store kv.Store) func(application *main.Application) {
	return func(a *main.Application) {
		a.Store = store
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

func (ts *testServer) put(t *testing.T, urlPath string, payload []byte) (int, http.Header, string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPut, ts.URL+urlPath, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
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

type simpleStore struct {
	store map[string]string
}

func (r *simpleStore) Set(key, value string) {
	r.store[key] = value
}
func (r *simpleStore) Get(key string) (value string, ok bool) {
	value, ok = r.store[key]
	return
}

type valuePayload struct {
	Value string
}

func valueToJSON(t *testing.T, value string) []byte {

	t.Helper()

	payload, err := json.Marshal(valuePayload{Value: value})

	if err != nil {
		t.Fatal(err)
	}

	return payload

}

func valueFromJSON(t *testing.T, data string) string {

	t.Helper()

	var payload valuePayload

	err := json.Unmarshal([]byte(data), &payload)

	if err != nil {
		t.Fatal(err)
	}

	return payload.Value

}
