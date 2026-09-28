package goastedit

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lordnynex/gest"
)

func TestReport(t *testing.T) {
	gest.Run(t, "Report", func(s *gest.S) {
		s.It("prints the comparison", func(t *gest.T) {
			text, err := comparison()
			t.Require(err).To(gest.BeNil())
			t.Expect(text).To(gest.Contain("replace body"))
			t.Expect(text).To(gest.Contain("search/replace"))
			t.Expect(text).To(gest.Contain("n/a"))
			fmt.Println(text)
		})
	})
}

func comparison() (string, error) {
	fixture, err := os.ReadFile("testdata/sample.go")
	if err != nil {
		return "", err
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.work")); err != nil {
		return "", fmt.Errorf("repo root: %w", err)
	}
	label, largest, err := LargestGoFile(root)
	if err != nil {
		return "", err
	}
	return Report(fixture, label, largest)
}
