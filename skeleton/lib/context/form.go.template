package context

import (
	"github.com/felipeweb/gopher-utils"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"reflect"
	"strings"
)

// Form representation
type Form interface {
	binding.Validator
}

func init() {
	binding.SetNameMapper(gopher_utils.ToSnakeCase)
}

// AssignForm assign form values back to the template data.
func AssignForm(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = gopher_utils.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

// GetSize get size validation
func GetSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

// Validate form
func Validate(errs binding.Errors, data map[string]interface{}, f Form, l macaron.Locale) binding.Errors {
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	AssignForm(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			name := field.Tag.Get("name")
			if len(name) == 0 {
				name = field.Name
			}

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = name + l.Tr("required")
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = name + l.Tr("dash")
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = name + l.Tr("dash")
			case binding.ERR_SIZE:
				data["ErrorMsg"] = name + l.Tr("must_size") + GetSize(field)
			default:
				data["ErrorMsg"] = l.Tr("unknown") + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}
