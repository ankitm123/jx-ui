package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"jx-ui/internal/server"
)

func TestPipelinesHandler(t *testing.T) {
	s := &server.Server{}
	s.Namespace = "jx-test"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pipelines", nil)
	w := httptest.NewRecorder()
	s.PipelinesHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ABC" {
		t.Errorf("expected ABC got %v", string(data))
	}
}
