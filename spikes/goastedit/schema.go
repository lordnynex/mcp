package goastedit

// Schema text is counted once per session. Search/replace is the baseline a
// client already pays. The splice and AST schemas are the incremental cost of
// adding those tools. They describe the payload shape, not every Go node kind.
const (
	SchemaSearch = `{"type":"object","description":"Replace an exact source span.","properties":{"old":{"type":"string"},"new":{"type":"string"}},"required":["old","new"]}`
	SchemaDiff   = `{"type":"object","description":"Apply a unified diff hunk.","properties":{"diff":{"type":"string"}},"required":["diff"]}`
	SchemaSplice = `{"type":"object","description":"Edit a Go declaration by name.","properties":{"op":{"enum":["replace_body","replace_decl","insert_before","insert_after","prepend_stmt","delete","add_import"]},"target":{"type":"string"},"content":{"type":"string"},"path":{"type":"string"}},"required":["op"]}`
	SchemaAST    = `{"type":"object","description":"Edit a Go declaration with a structural syntax node.","properties":{"op":{"type":"string"},"target":{"type":"string"},"path":{"type":"string"},"node":{"type":"object","description":"Node with a kind field and child nodes. Positions and comments are omitted.","additionalProperties":true}},"required":["op"]}`
	SchemaIntent = `{"type":"object","description":"Apply a refactor that needs no replacement source.","properties":{"op":{"enum":["delete","add_import"]},"target":{"type":"string"},"path":{"type":"string"}},"required":["op"]}`
)
