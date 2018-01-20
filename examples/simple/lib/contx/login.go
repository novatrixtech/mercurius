package contx

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type Login struct {
	Username string `binding:"Required"`
	Password string `binding:"Required"`
}

func (f *Login) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return Validate(errs, ctx.Data, f, ctx.Locale)
}
