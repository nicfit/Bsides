package m3u

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding"
)

func TestLocalEncoding(t *testing.T) {
	// Happy path
	lc, err := getLocalEncoding()
	assert.NoError(t, err)
	assert.NotEmpty(t, lc)

	// Not found
	LC_ALL := os.Getenv("LC_ALL")
	LC_CTYPE := os.Getenv("LC_CTYPE")
	LANG := os.Getenv("LANG")
	defer func() {
		os.Setenv("LC_ALL", LC_ALL)
		os.Setenv("LC_CTYPE", LC_CTYPE)
		os.Setenv("LANG", LANG)
	}()
	assert.NoError(t, os.Setenv("LC_ALL", ""))
	assert.NoError(t, os.Setenv("LC_CTYPE", ""))
	assert.NoError(t, os.Setenv("LANG", ""))

	lc, err = getLocalEncoding()
	assert.NoError(t, err)
	switch runtime.GOOS {
	case "windows":
		assert.NotEmpty(t, lc)
	default:
		assert.Equal(t, encoding.Nop, lc)
	}
}
