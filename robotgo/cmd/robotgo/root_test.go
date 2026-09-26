package main

import (
	"testing"

	"github.com/lordnynex/gest"
)

func TestRoot(t *testing.T) {
	gest.Run(t, "Root", func(s *gest.S) {
		s.It("is a cobra command named robotgo", func(t *gest.T) {
			t.Expect(rootCmd).NotTo(gest.BeNil())
			t.Expect(rootCmd.Use).To(gest.Equal("robotgo"))
		})
	})
}
