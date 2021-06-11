package query

import (
	"log"
	"testing"
	"unicode"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestFuzzQueryNonASCII(t *testing.T) {
	t.Parallel()
	f := fuzz.New()
	var query string
	for i := 0; i < 100000; {
		f.Fuzz(&query)
		if !isASCII([]byte(query)) {
			continue
		}
		require.NotPanics(t, func() {
			Parse(query)
		}, "panicked with input %s", string(query))
		i++
	}
}

func TestFuzzQueryASCII(t *testing.T) {
	t.Parallel()
	f := fuzz.New()
	var query []byte
	for i := 0; i < 100000; {
		f.Fuzz(&query)
		if isASCII(query) {
			continue
		}
		require.NotPanics(t, func() {
			Parse(string(query))
		}, "panicked with input %s", string(query))
		i++
	}
}

func TestFuzzRegressions(t *testing.T) {
	crashers := []string{`query($~\344\334\234\344\334\344\234�d44\201"`}
	for _, crash := range crashers {
		require.NotPanics(t, func() {
			_, err := Parse(crash)
			if err == nil {
				log.Fatalf("found a regression with %s", crash)
			}
		}, "panicked for query: %s", crash)
	}
}

func isASCII(b []byte) bool {
	for i := range b {
		if b[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
