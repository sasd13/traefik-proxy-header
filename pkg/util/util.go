//nolint:revive
package util

import (
	"net/http"
)

// CopyHeader copy one header to another.
func CopyHeader(r *http.Request, orig, dest string, overwrite bool) {
	value := r.Header.Get(orig)
	if value != "" {
		existing := r.Header.Get(dest)
		if existing == "" || overwrite {
			r.Header.Set(dest, value)
		}
	}
}
