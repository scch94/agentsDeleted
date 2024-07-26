package modelUtils

import "strings"

type Table struct {
	TableName     string
	Conditional   string
	QueryToDelete strings.Builder
}
