package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const origin = "localhost"

func TestCorsHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodOptions, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Origin", origin)
	r.Header.Set("Access-Control-Request-Method", http.MethodGet)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app := &application{
		config: config{
			cors: struct{ trustedOrigins []string }{trustedOrigins: []string{origin}},
		},
	}

	app.enableCORS(next).ServeHTTP(rr, r)

	rs := rr.Result()

	varyOptions := rs.Header.Get("Vary")
	if varyOptions != "Origin" {
		t.Errorf("want %q; got %q", "Origin", varyOptions)
	}

	originOptions := rs.Header.Get("Access-Control-Allow-Origin")
	if originOptions != origin {
		t.Errorf("want %q; got %q", origin, originOptions)
	}

	methodOptions := rs.Header.Get("Access-Control-Allow-Methods")
	if methodOptions != "OPTIONS, PUT, PATCH, DELETE" {
		t.Errorf("want %q; got %q", "OPTIONS, PUT, PATCH, DELETE", methodOptions)
	}

	headersOptions := rs.Header.Get("Access-Control-Allow-Headers")
	if headersOptions != "Authorization, Content-Type" {
		t.Errorf("want %q; got %q", "Authorization, Content-Type", headersOptions)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
}
