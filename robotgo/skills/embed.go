package skills

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed observe-then-act/SKILL.md keyboard/SKILL.md mouse/SKILL.md screenshot/SKILL.md aim-center/SKILL.md
var FS embed.FS

const (
	IndexURI    = "robotgo://skills"
	URIPrefix   = "robotgo://skills/"
	URITemplate = "robotgo://skills/{name}"
)

// Names are skill directory names.
var Names = []string{"aim-center", "observe-then-act", "keyboard", "mouse", "screenshot"}

// Read returns the SKILL.md body for name.
func Read(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("skill name is required")
	}
	b, err := FS.ReadFile(name + "/SKILL.md")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Index lists skill names.
func Index() []string {
	var out []string
	_ = fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if strings.HasSuffix(path, "/SKILL.md") {
			out = append(out, strings.TrimSuffix(path, "/SKILL.md"))
		}
		return nil
	})
	return out
}
