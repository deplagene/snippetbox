package snippetbox

import "embed"

//go:embed sql/schema/*.sql
var SchemaFS embed.FS
