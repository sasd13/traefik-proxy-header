// Package traefik_proxy_header provides functionality for adapting request headers.
// revive:disable-next-line var-naming.
package traefik_proxy_header

import (
	"context"
	"log"
	"net/http"

	adp "github.com/sasd13/traefik-proxy-header/pkg/adapter"
)

// Config the plugin configuration.
type Config struct {
	Debug    bool     `json:"debug,omitempty"`
	Adapters []string `json:"adapters,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Debug:    false,
		Adapters: []string{},
	}
}

// ProxyHeader plugin struct.
type ProxyHeader struct {
	name     string
	next     http.Handler
	debug    bool
	adapters []string
}

// New created a new ProxyHeader plugin.
// revive:disable-next-line unused-parameter.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &ProxyHeader{
		name:     name,
		next:     next,
		debug:    config.Debug,
		adapters: config.Adapters,
	}, nil
}

func (p *ProxyHeader) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s\n", r.Method, r.URL)

	if p.debug {
		log.Println("Headers:", r.Header)
	}

	for _, typ := range p.adapters {
		adapter := adp.AdapterFactory(typ)
		if adapter != nil {
			adapter.Headers(r)
		} else {
			log.Printf("Unknown adapter type: %s\n", typ)
		}
	}

	p.next.ServeHTTP(rw, r)
}
