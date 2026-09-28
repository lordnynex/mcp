package goastedit

import (
	"fmt"
	"strings"
)

// unifiedDiff is a unified diff of before and after with three lines of context.
// It is the hunk a model would emit for the minimal byte edit, before gofmt.
func unifiedDiff(before, after string) string {
	ops := lineEdits(splitLines(before), splitLines(after))
	const ctx = 3
	var b strings.Builder
	i := 0
	for i < len(ops) {
		if ops[i].kind == opEqual {
			i++
			continue
		}
		end := i
		for end < len(ops) {
			if ops[end].kind != opEqual {
				end++
				continue
			}
			k := end
			for k < len(ops) && ops[k].kind == opEqual {
				k++
			}
			if k == len(ops) || k-end > 2*ctx {
				break
			}
			end = k
		}
		pre := contextStart(ops, i, ctx)
		post := end
		got := 0
		for post < len(ops) && got < ctx && ops[post].kind == opEqual {
			post++
			got++
		}
		aLine, bLine := 1, 1
		for k := 0; k < pre; k++ {
			switch ops[k].kind {
			case opEqual:
				aLine++
				bLine++
			case opDel:
				aLine++
			case opIns:
				bLine++
			}
		}
		aCount, bCount := 0, 0
		var body strings.Builder
		for k := pre; k < post; k++ {
			switch ops[k].kind {
			case opEqual:
				fmt.Fprintf(&body, " %s\n", ops[k].line)
				aCount++
				bCount++
			case opDel:
				fmt.Fprintf(&body, "-%s\n", ops[k].line)
				aCount++
			case opIns:
				fmt.Fprintf(&body, "+%s\n", ops[k].line)
				bCount++
			}
		}
		fmt.Fprintf(&b, "@@ -%d,%d +%d,%d @@\n%s", aLine, aCount, bLine, bCount, body.String())
		i = end
	}
	return b.String()
}

func contextStart(ops []lineOp, change, ctx int) int {
	pre := change
	got := 0
	for pre > 0 && got < ctx && ops[pre-1].kind == opEqual {
		pre--
		got++
	}
	return pre
}

type opKind int

const (
	opEqual opKind = iota
	opDel
	opIns
)

type lineOp struct {
	kind opKind
	line string
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	if strings.HasSuffix(s, "\n") {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func lineEdits(a, b []string) []lineOp {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}
	out := make([]lineOp, 0, n+m)
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			out = append(out, lineOp{opEqual, a[i]})
			i++
			j++
			continue
		}
		if dp[i+1][j] >= dp[i][j+1] {
			out = append(out, lineOp{opDel, a[i]})
			i++
			continue
		}
		out = append(out, lineOp{opIns, b[j]})
		j++
	}
	for i < n {
		out = append(out, lineOp{opDel, a[i]})
		i++
	}
	for j < m {
		out = append(out, lineOp{opIns, b[j]})
		j++
	}
	return out
}
