// Package traefik_proxy_header_test provides the tests.
// revive:disable-next-line var-naming.
package traefik_proxy_header_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	plugin "github.com/sasd13/traefik-proxy-header"
	"github.com/stretchr/testify/assert"
)

func TestAdapter(t *testing.T) {
	cfg := plugin.CreateConfig()
	cfg.Debug = true
	cfg.Adapters = []string{"cloudflare"}

	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	handler, err := plugin.New(ctx, next, cfg, "proxy-header-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	req.Header.Set("Location", "https://api.nasa.gov/planetary/apod")
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Result().StatusCode)
}
