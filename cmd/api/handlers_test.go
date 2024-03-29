package main

import (
	"bytes"
	"encoding/json"
	"github.com/alisher0594/reddit-old/internal/data/entitys"
	"net/http"
	"reflect"
	"testing"
)

func TestHandlers(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name     string
		testFunc func(t *testing.T, ts *testServer)
	}{
		{
			name:     "HealthcheckHandler",
			testFunc: HealthcheckHandler,
		},
		{
			name:     "CreatePostHandler",
			testFunc: CreatePostHandler,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, ts)
		})
	}
}

func HealthcheckHandler(t *testing.T, ts *testServer) {
	code, _, body := ts.get(t, "/healthcheck")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	want := []byte(`{
	"status": "available",
	"system_info": {
		"environment": "staging"
	}
}
`)
	if !bytes.Contains(body, want) {
		t.Errorf("want body to equal %q, got: %q", string(want), string(body))
	}
}

func CreatePostHandler(t *testing.T, ts *testServer) {
	code, _, body := ts.get(t, "/v1/posts/13")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	var result struct {
		Post *entitys.Post `json:"post"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("want %d; got %d", http.StatusOK, code)
	}

	want := &entitys.Post{
		ID:       13,
		Title:    "Mocked post",
		Author:   "t2_fed2ere13",
		Link:     "https://old.reddit.com/user/fed2ere13",
		Content:  "content of mocked post",
		Score:    100500,
		Promoted: false,
		NSFW:     false,
		Version:  1,
	}

	if !reflect.DeepEqual(want, result.Post) {
		t.Errorf("\nwant: %v;\ngot %v", want, result)
	}
}
