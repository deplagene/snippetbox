package snippetbox

import "embed"

//go:embed all:sql/schema
var SchemaFS embed.FS

//go:embed all:ui
var UIFS embed.FS