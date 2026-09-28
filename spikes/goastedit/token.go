package goastedit

import (
	"fmt"
	"sync"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

// Usage is the size of one payload.
type Usage struct {
	Bytes  int
	Tokens int
}

var (
	encOnce sync.Once
	enc     *tiktoken.Tiktoken
	encErr  error
)

// Measure counts UTF-8 bytes and o200k_base tokens. Tokens estimate model
// usage; they are not a bill from a live call.
func Measure(s string) (Usage, error) {
	encOnce.Do(func() {
		enc, encErr = tiktoken.GetEncoding("o200k_base")
	})
	if encErr != nil {
		return Usage{}, fmt.Errorf("o200k_base: %w", encErr)
	}
	return Usage{Bytes: len(s), Tokens: len(enc.Encode(s, nil, nil))}, nil
}
