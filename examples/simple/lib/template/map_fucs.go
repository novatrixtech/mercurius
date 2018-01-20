package template

import (
	"html/template"

	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
)

func FuncMaps() []template.FuncMap {
	return []template.FuncMap{
		map[string]interface{}{
			"Tr": contx.I18n,
		}}
}
