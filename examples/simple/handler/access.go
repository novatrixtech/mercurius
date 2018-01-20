package handler

import (
	"net/http"

	"log"

	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
	"github.com/novatrixtech/mercurius/examples/simple/lib/query"
	"github.com/novatrixtech/mercurius/examples/simple/model"
	"github.com/novatrixtech/mercurius/examples/simple/repo"
)

func ListAccessPage(ctx *contx.Context) {
	ctx.Data["rows"] = 0
	ctx.HTML(http.StatusOK, "list")
}

func ListAccessBy(ctx *contx.Context) {
	access, err := list(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	ctx.Data["access"] = access
	ctx.Data["rows"] = len(access)
	ctx.HTML(http.StatusOK, "list")
}

func ListAccessForApi(ctx *contx.Context) {
	access, err := list(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, access)
}

func list(ctx *contx.Context) ([]model.Access, error) {
	fields := make(map[string]string)
	if marca := ctx.Query("marca"); marca != "" {
		fields["marca"] = marca
	}
	if modelo := ctx.Query("modelo"); modelo != "" {
		fields["modelo"] = modelo
	}
	if combustivel := ctx.Query("combustivel"); combustivel != "" {
		fields["combustivel"] = combustivel
	}
	if propriedade := ctx.Query("propriedade"); propriedade != "" {
		fields["propriedade"] = propriedade
	}
	if ano := ctx.Query("ano"); ano != "" {
		fields["ano"] = ano
	}
	if dataInicio := ctx.Query("dataInicio"); dataInicio != "" {
		fields["dataInicio"] = dataInicio

	}
	if dataFim := ctx.Query("dataFim"); dataFim != "" {
		fields["dataFim"] = dataFim

	}

	objRepo, err := repo.NewAccessRepository()
	if err != nil {
		return nil, err
	}
	access, err := objRepo.FindAllBy(query.Build(fields), ctx.Cache)
	if err != nil {
		return nil, err
	}
	return access, nil
}
