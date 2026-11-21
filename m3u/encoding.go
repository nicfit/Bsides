//go:build linux || darwin

package m3u

import (
	"os"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

func getLocalEncoding() (encoding.Encoding, error) {
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if val := os.Getenv(env); val != "" {
			if _, charset, found := strings.Cut(val, "."); found {
				// Strip modifiers like @euro if present
				charset, _, _ = strings.Cut(charset, "@")
				return ianaindex.IANA.Encoding(charset)
			}
		}
	}
	return encoding.Nop, nil // Default to raw/UTF-8 if unset
}
