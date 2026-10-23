//go:build utest

package globalhelper

import (
	"testing"
)

// FuzzConvertSpecNameToFileName checks that any spec name converts to a string
// that is safe to use as a report file name: no path separators, no spaces or
// commas, and no uppercase ASCII letters.
func FuzzConvertSpecNameToFileName(f *testing.F) {
	f.Add("should block ingress with TLS")
	f.Add("Affiliated Certification, operator is certified")
	f.Add("a/b/c d-e, f")
	f.Add("İstanbul 5G")
	f.Add("")

	f.Fuzz(func(t *testing.T, specName string) {
		fileName := ConvertSpecNameToFileName(specName)

		for _, r := range fileName {
			if r == ' ' || r == ',' || r == '/' || r == '-' {
				t.Fatalf("unexpected %q in %q (from %q)", r, fileName, specName)
			}

			if r >= 'A' && r <= 'Z' {
				t.Fatalf("unexpected uppercase %q in %q (from %q)", r, fileName, specName)
			}
		}
	})
}
