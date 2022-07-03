package integration_test

import (
	. "github.com/Eun/go-hit"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	attempts   = 20
	host       = "app:4000"
	healthPath = "http://" + host + "/healthcheck"

	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

func TestHTTPCreatePost(t *testing.T) {
	body := `{
		"author": "t2_ad4few0q",
		"link": "https://google.com",
		"title": "Test HTTP Create Post",
		"content": "Some content",
		"promoted": false,
		"nsfw": false
	}`
	Test(t,
		Description("CreatePost Success"),
		Post(basePath+"/posts"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),

		Expect().Status().Equal(http.StatusCreated),
		Expect().Headers("Location").Contains("/v1/post/49"),
		Expect().Body().JSON().JQ(".envelope.post.author").Equal("t2_ad4few0q"),
		Expect().Body().JSON().JQ(".envelope.post.link").Equal("https://google.com"),
		Expect().Body().JSON().JQ(".envelope.post.title").Equal("Test HTTP Create Post"),
		Expect().Body().JSON().JQ(".envelope.post.content").Equal("Some content"),
		Expect().Body().JSON().JQ(".envelope.post.promoted").Equal(false),
		Expect().Body().JSON().JQ(".envelope.post.nsfw").Equal(true),
	)

	body = `{
		"author": "t9_ad4few0q",
		"link": "https://google.com",
		"title": "Test HTTP Create Post",
		"content": "Some content",
		"promoted": false,
		"nsfw": false
	}`

	Test(t,
		Description("CreatePost Fail"),
		Post(basePath+"/posts"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),

		Expect().Status().Equal(http.StatusUnprocessableEntity),
		Expect().Body().JSON().JQ(".error.author").Equal("must have prefix: t2_"),
	)
}
