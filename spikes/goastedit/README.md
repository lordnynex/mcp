# goastedit

Measurement spike for Go edits. It compares the token cost of a text patch, a symbol-addressed splice, structural `go/ast` JSON, and an intent call that names a refactor without sending source. It is not an MCP server.

`spikes/` is the parent directory so a later TypeScript spike can sit beside this one.

## Run

From this directory:

```bash
go test ./...
go test -run TestReport -v
```

`go test ./...` from the repo root tests the root module only. This module is listed in `go.work` so it builds with the workspace toolchain.

## What is counted

Tokens are `o200k_base` via `github.com/pkoukk/tiktoken-go`, plus raw UTF-8 bytes. They estimate model usage. Nothing here calls a model.

Each scenario writes the new Go once. The arms are different renderings of that edit:

- **search/replace** quotes an old span and the new text (`SEARCH` / `REPLACE`). A replacement or delete quotes that span. An insertion quotes the enclosing body, declaration, or import block, because an empty old string is not a usable match.
- **unified diff** is a hunk with three lines of context around the minimal byte edit, before `gofmt`.
- **symbol splice** is JSON: the operation, the declaration address (`(*Box).Get`, `Other.Label`, `Helper`), and the new Go source only.
- **ast json** is the same operation with a structural node (`{"kind":"AssignStmt",...}`). Byte positions and comments are omitted.
- **intent** is only for delete and add-import, where the model sends no replacement source. Other scenarios are `n/a`.

Search/replace and unified diff are raw text. The other arms are JSON. A transport that JSON-encodes every argument would add escaping to the text arms.

Schema strings are counted once per session. Search/replace and unified diff are treated as tools a client already pays for. The splice, AST, and intent schemas are the incremental cost of adding them. The AST schema is a generic node object, not a closed listing of every `go/ast` type.

Tool results are what a server would send back after `gofmt`, independent of which arm asked for the edit: a one-line summary, a unified diff, the full file, or structural JSON of the whole file.

The read section compares source, a signature outline, and structural JSON. The second file is the largest parseable `.go` file in the repo, skipping hidden directories and `vendor`.

## Results

`go test -run TestReport -v` on this tree:

### Schema (per session)

| Arm            | Bytes | Tokens | Role                  |
| -------------- | ----: | -----: | --------------------- |
| search/replace |   151 |     36 | baseline already paid |
| unified diff   |   120 |     28 | baseline already paid |
| symbol splice  |   291 |     67 | incremental           |
| ast json       |   339 |     74 | incremental           |
| intent         |   207 |     48 | incremental           |

### Read

| File                   | View     | Bytes | Tokens |
| ---------------------- | -------- | ----: | -----: |
| sample.go              | source   |  2086 |    524 |
| sample.go              | outline  |   280 |     85 |
| sample.go              | ast json | 11404 |   3305 |
| robotgo/tools/batch.go | source   | 13162 |   3706 |
| robotgo/tools/batch.go | outline  |   667 |    176 |
| robotgo/tools/batch.go | ast json | 70335 |  20020 |

### Model output

| Scenario          | Arm            | Bytes | Tokens |
| ----------------- | -------------- | ----: | -----: |
| replace body      | search/replace |   333 |    102 |
| replace body      | unified diff   |   570 |    173 |
| replace body      | symbol splice  |   165 |     53 |
| replace body      | ast json       |   792 |    232 |
| replace body      | intent         |   n/a |    n/a |
| add method        | search/replace |   638 |    180 |
| add method        | unified diff   |   316 |     94 |
| add method        | symbol splice  |   173 |     52 |
| add method        | ast json       |   881 |    257 |
| add method        | intent         |   n/a |    n/a |
| change signature  | search/replace |   529 |    135 |
| change signature  | unified diff   |   616 |    169 |
| change signature  | symbol splice  |   291 |     79 |
| change signature  | ast json       |  1321 |    375 |
| change signature  | intent         |   n/a |    n/a |
| delete helper     | search/replace |   309 |     82 |
| delete helper     | unified diff   |   487 |    137 |
| delete helper     | symbol splice  |    33 |      9 |
| delete helper     | ast json       |    33 |      9 |
| delete helper     | intent         |    33 |      9 |
| add import        | search/replace |    83 |     26 |
| add import        | unified diff   |   119 |     39 |
| add import        | symbol splice  |    36 |     10 |
| add import        | ast json       |   129 |     35 |
| add import        | intent         |    36 |     10 |
| prepend statement | search/replace |   379 |    108 |
| prepend statement | unified diff   |   224 |     71 |
| prepend statement | symbol splice  |   102 |     30 |
| prepend statement | ast json       |   353 |    100 |
| prepend statement | intent         |   n/a |    n/a |

### Tool result

| Scenario          | Result   | Bytes | Tokens |
| ----------------- | -------- | ----: | -----: |
| replace body      | summary  |    27 |      8 |
| replace body      | diff     |   570 |    173 |
| replace body      | file     |  1955 |    479 |
| replace body      | ast json | 10118 |   2923 |
| add method        | summary  |    25 |      7 |
| add method        | diff     |   314 |     93 |
| add method        | file     |  2199 |    555 |
| add method        | ast json | 12235 |   3547 |
| change signature  | summary  |    27 |      4 |
| change signature  | diff     |   616 |    169 |
| change signature  | file     |  2013 |    500 |
| change signature  | ast json | 10572 |   3059 |
| delete helper     | summary  |    14 |      2 |
| delete helper     | diff     |   485 |    136 |
| delete helper     | file     |  1791 |    448 |
| delete helper     | ast json |  9297 |   2696 |
| add import        | summary  |    20 |      3 |
| add import        | diff     |    65 |     28 |
| add import        | file     |  2097 |    527 |
| add import        | ast json | 11490 |   3328 |
| prepend statement | summary  |    33 |      8 |
| prepend statement | diff     |   227 |     73 |
| prepend statement | file     |  2132 |    537 |
| prepend statement | ast json | 11707 |   3390 |

## What this shows

Writing `go/ast` nodes does not conserve tokens when the edit contains new code. Replacing a method body costs 232 tokens as a tree, 53 as a symbol splice, and 102 as search/replace. Changing a signature is 375 versus 79 and 135. The tree is several times the Go it describes, even with positions and comments left out. A comment-preserving encoding would be larger. The AST schema counted above is only 74 tokens because it does not enumerate node kinds.

Naming a declaration and sending Go source is the smallest model output on every scenario here, including a full body replacement, not only delete, add-import, and prepend. The saving is that the model does not repeat the old span. Delete and add-import need no new source at all: splice and intent are 9–10 tokens, against 26–137 for text.

That saving is smaller than the cost of reading the result back as a file (448–555 tokens) or as a tree (2,696–3,547 on the fixture, and 20,020 versus 3,706 source tokens on `robotgo/tools/batch.go`). A one-line summary is 2–8 tokens. An outline is the cheap way to read: 85 tokens versus 524 for the fixture, and 176 versus 3,706 for `batch.go`.
