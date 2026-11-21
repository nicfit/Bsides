//go:build windows

package m3u

import (
	"fmt"

	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	getACP   = kernel32.NewProc("GetACP")
)

func getLocalEncoding() (encoding.Encoding, error) {
	ret, _, _ := getACP.Call()
	codePage := fmt.Sprintf("windows-%d", ret)

	enc, err := ianaindex.IANA.Encoding(codePage)
	if err != nil || enc == nil {
		// Fallback for code pages not named "windows-XXX" (e.g., cp932 -> Shift_JIS)
		return ianaindex.IANA.Encoding(fmt.Sprintf("cp%d", ret))
	}
	return enc, nil
}
