package context

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

// Login form representation
type Login struct {
	Username string `binding:"Required"`
	Password string `binding:"Required"`
}

// Validate Login form
func (f *Login) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return Validate(errs, ctx.Data, f, ctx.Locale)
}
