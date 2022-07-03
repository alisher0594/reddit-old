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
	healthPath = "http://" + host + "/v1/healthcheck"

	basePath = "http://" + host + "/v1"
)

// RabbitMQ RPC Client: getHistory.
func TestMock(t *testing.T) {

}

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
