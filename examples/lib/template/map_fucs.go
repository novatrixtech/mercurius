package template

import (
	"html/template"

	"github.com/novatrixtech/mercurius/examples/lib/context"
)

func FuncMaps() []template.FuncMap {
	return []template.FuncMap{
		map[string]interface{}{
			"Tr": context.I18n,
		}}
}
