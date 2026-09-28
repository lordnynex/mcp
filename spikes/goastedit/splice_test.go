package goastedit

import (
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/lordnynex/gest"
)

func TestSplice(t *testing.T) {
	gest.Run(t, "Splice", func(s *gest.S) {
		s.It("keeps the doc comment and signature when replacing a body", func(t *gest.T) {
			out := applyNamed(t, "replace body")
			text := string(out)
			t.Expect(text).To(gest.Contain("Get returns the value stored for key."))
			t.Expect(text).To(gest.Contain("func (b *Box) Get(key string) string"))
			t.Expect(text).To(gest.Contain(`return "missing"`))
			t.Expect(text).NotTo(gest.Contain("b.hits++"))
		})

		s.It("removes the doc comment and keeps the neighboring comment on delete", func(t *gest.T) {
			out := applyNamed(t, "delete helper")
			text := string(out)
			t.Expect(text).NotTo(gest.Contain("func Helper"))
			t.Expect(text).NotTo(gest.Contain("Helper formats a label."))
			t.Expect(text).To(gest.Contain("must survive deletion"))
			t.Expect(text).To(gest.Contain("func Neighbor() string"))
		})

		s.It("joins an existing import group", func(t *gest.T) {
			out := applyNamed(t, "add import")
			text := string(out)
			t.Expect(text).To(gest.Contain("import ("))
			t.Expect(text).To(gest.Contain(`"fmt"`))
			t.Expect(text).To(gest.Contain(`"strconv"`))
			t.Expect(text).To(gest.Contain(`"strings"`))
		})

		s.It("lists both matches for an ambiguous name", func(t *gest.T) {
			idx, err := Parse(readFixture(t))
			t.Require(err).To(gest.BeNil())
			_, err = idx.Lookup("Get")
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(err.Error()).To(gest.Contain("(*Box).Get"))
			t.Expect(err.Error()).To(gest.Contain("(*Other).Get"))
		})

		s.It("prepends a statement without dropping the rest of the body", func(t *gest.T) {
			out := applyNamed(t, "prepend statement")
			text := string(out)
			t.Expect(text).To(gest.Contain(`strings.Contains(key, " ")`))
			t.Expect(text).To(gest.Contain("b.items == nil"))
			t.Expect(text).To(gest.Contain("b.items[key] = strings.TrimSpace(val)"))
		})

		s.It("resolves a value receiver as T.Method", func(t *gest.T) {
			idx, err := Parse(readFixture(t))
			t.Require(err).To(gest.BeNil())
			got, err := idx.Lookup("Other.Label")
			t.Require(err).To(gest.BeNil())
			t.Expect(got.Addr).To(gest.Equal("Other.Label"))
			_, err = idx.Lookup("(*Other).Label")
			t.Expect(err).NotTo(gest.BeNil())
		})

		s.It("formats every scenario into parseable Go", func(t *gest.T) {
			src := readFixture(t)
			for _, sc := range Scenarios() {
				out, err := Apply(src, sc.Req)
				t.Require(err).To(gest.BeNil())
				again, err := format.Source(out)
				t.Require(err).To(gest.BeNil())
				t.Expect(string(again)).To(gest.Equal(string(out)))
				_, err = parser.ParseFile(token.NewFileSet(), sc.Name+".go", out, parser.SkipObjectResolution)
				t.Require(err).To(gest.BeNil())
			}
		})
	})
}

func TestPayloads(t *testing.T) {
	gest.Run(t, "Payloads", func(s *gest.S) {
		s.It("quotes the old body for search/replace and not for a prepend splice", func(t *gest.T) {
			arms := armsNamed(t, "prepend statement")
			search := armNamed(t, arms, "search/replace")
			splice := armNamed(t, arms, "symbol splice")
			t.Expect(search.Payload).To(gest.Contain("b.items == nil"))
			t.Expect(splice.Payload).NotTo(gest.Contain("b.items == nil"))
			t.Expect(splice.Payload).To(gest.Contain("strings.Contains"))
		})

		s.It("has no intent payload for a body replacement", func(t *gest.T) {
			intent := armNamed(t, armsNamed(t, "replace body"), "intent")
			t.Expect(intent.NA).To(gest.Equal(true))
		})

		s.It("names the helper in a delete intent without quoting its body", func(t *gest.T) {
			intent := armNamed(t, armsNamed(t, "delete helper"), "intent")
			t.Expect(intent.NA).To(gest.Equal(false))
			t.Expect(intent.Payload).To(gest.Contain("delete"))
			t.Expect(intent.Payload).To(gest.Contain("Helper"))
			t.Expect(intent.Payload).NotTo(gest.Contain("func Helper"))
		})
	})
}

func readFixture(t *gest.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/sample.go")
	t.Require(err).To(gest.BeNil())
	return b
}

func applyNamed(t *gest.T, name string) []byte {
	t.Helper()
	out, err := Apply(readFixture(t), scenarioNamed(t, name).Req)
	t.Require(err).To(gest.BeNil())
	return out
}

func armsNamed(t *gest.T, name string) []Arm {
	t.Helper()
	arms, err := Payloads(readFixture(t), scenarioNamed(t, name).Req)
	t.Require(err).To(gest.BeNil())
	return arms
}

func scenarioNamed(t *gest.T, name string) Scenario {
	t.Helper()
	for _, sc := range Scenarios() {
		if sc.Name == name {
			return sc
		}
	}
	t.Require(name).To(gest.Equal(""))
	return Scenario{}
}

func armNamed(t *gest.T, arms []Arm, name string) Arm {
	t.Helper()
	for _, arm := range arms {
		if arm.Name == name {
			return arm
		}
	}
	t.Require(name).To(gest.Equal(""))
	return Arm{}
}
