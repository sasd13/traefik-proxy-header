//nolint:revive
package adapter

import (
	"net/http"
)

type Adapter interface {
	Headers(r *http.Request)
}

func AdapterFactory(typ string) Adapter { //nolint:ireturn
	switch typ {
	case "cloudflare":
		return CloudflareAdapter{}
	default:
		return nil
	}
}
