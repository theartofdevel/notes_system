package metric

import (
	"net/http"
)

const (
	HEARTBEAT_URL = "/api/heartbeat"
	TEST_URL      = "/api/test"
)

func Heartbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}

func Test(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}
