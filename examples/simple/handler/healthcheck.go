package handler

import (
	"net/http"
	"time"

	"github.com/novatrixtech/mercurius/examples/simple/lib/context"
)

/*
HealthCheck - Check if the application is up and working
*/
func HealthCheck(ctx *context.Context) {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	ctx.Data["dataehora"] = t.String()
	ctx.HTML(http.StatusOK, "index")
}
