package cors

import (
	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
	"gopkg.in/macaron.v1"
)

func Cors() macaron.Handler {
	return func(ctx *contx.Context) {
		ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		ctx.Resp.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Resp.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Next()
	}
}
