//nolint:revive
package adapter

import (
	"net/http"
	"strings"

	"github.com/sasd13/traefik-proxy-header/pkg/util"
)

const (
	X_CLIENT_IP     = "X-Client-IP"
	X_FORWARDED_FOR = "X-Forwarded-For"
)

type CloudflareAdapter struct{}

func (a CloudflareAdapter) Headers(r *http.Request) {
	addCustomHeaders(r)
	refineProxyHeaders(r)
}

func addCustomHeaders(r *http.Request) {
	util.CopyHeader(r, "CF-Connecting-IP", X_CLIENT_IP, false)
	util.CopyHeader(r, "CF-Connecting-IPv6", "X-Client-IPv6", false)
	util.CopyHeader(r, "CF-IPCountry", "X-Client-Country", false)
}

func refineProxyHeaders(r *http.Request) {
	xClientIP := r.Header.Get(X_CLIENT_IP)
	if xClientIP == "" {
		return
	}

	xForwardedFor := r.Header.Get(X_FORWARDED_FOR)
	if xForwardedFor == "" {
		// Set it directly if empty
		r.Header.Set(X_FORWARDED_FOR, xClientIP)
		return
	}

	entries := strings.Split(xForwardedFor, ",")
	for i := range entries {
		entries[i] = strings.TrimSpace(entries[i])
	}

	// If X-Client-IP is already at the beginning, do nothing
	if len(entries) > 0 && entries[0] == xClientIP {
		return
	}

	// Remove X-Client-IP if it's anywhere else in the list to avoid duplicates
	filtered := []string{}
	for _, entry := range entries {
		if entry != xClientIP {
			filtered = append(filtered, entry)
		}
	}

	// Prepend X-Client-IP
	newForwardedFor := xClientIP
	if len(filtered) > 0 {
		newForwardedFor += ", " + strings.Join(filtered, ", ")
	}

	r.Header.Set(X_FORWARDED_FOR, newForwardedFor)
}
