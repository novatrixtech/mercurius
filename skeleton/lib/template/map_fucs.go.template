package template

import (
	"html/template"

	"{{.AppPath}}/lib/context"
)

func FuncMaps() []template.FuncMap {
	return []template.FuncMap{
		map[string]interface{}{
			"Tr": context.I18n,
		}}
}
