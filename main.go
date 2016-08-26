package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/flosch/pongo2-addons"
	"github.com/iris-contrib/template/django"
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/xgfone/go-utils/log"
	"github.com/xgfone/gobtweb/g"
)

func handlerError(ctx *iris.Context, err error) {
	if err != nil {
		log.Errorj("Failed to return the response", "err", err)
		ctx.EmitError(iris.StatusInternalServerError)
	}
}

func index(ctx *iris.Context) {
	err := ctx.Render("index.html", nil)
	handlerError(ctx, err)
}

// func downloadTorrent(ctx *iris.Context) {
// 	infohash := ctx.URLParam("torrent")
// 	if infohash == "" {
// 		ctx.EmitError(iris.StatusNotFound)
// 		return
// 	} else if len(infohash) != 40 {
// 		ctx.EmitError(iris.StatusBadRequest)
// 		return
// 	}

// 	//data := map[string]interface{}{}
// 	data, err := GetTorrentByInfohashFromSE(g.Conf.ElasticClient, infohash)
// 	if err != nil {
// 		handlerError(err)
// 		return
// 	}

// 	err = ctx.Data(iris.StatusOK, data)
// 	handlerError(ctx, err)
// }

func search(ctx *iris.Context) {
	query := strings.TrimSpace(ctx.URLParam("q"))
	if query == "" {
		query = strings.TrimSpace(ctx.FormValueString("key"))
		if query == "" {
			ctx.Redirect("/")
			return
		}
	}

	page, err := ctx.URLParamInt("page")
	if err != nil || page < 1 {
		page = 1
	}

	size := g.Conf.PageSize
	from := size * (page - 1)

	data, total, err := SearchKeyword(g.ElasticClient, g.Repository, query, from, size)
	if err != nil {
		handlerError(ctx, err)
		return
	} else if data == nil || len(data) == 0 {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	urlPattern := fmt.Sprintf("/search?key=%v&page=", query)
	pn := NewPagination(total, page, size, urlPattern)
	pn.SetNumLinks(10).Init()
	err = ctx.Render("search2.html", map[string]interface{}{
		"query":      query,
		"torrents":   data,
		"total":      total,
		"pagination": *pn,
	})
	handlerError(ctx, err)
}

func main() {
	g.Init(os.Args[1])

	cfg := config.Default()
	cfg.ProfilePath = "/debug"

	web := iris.New(cfg)

	templateDir := g.Conf.TemplateDirectory
	django_cfg := django.DefaultConfig()
	django_cfg.DebugTemplates = g.Conf.TemplateDebug
	web.UseTemplate(django.New(django_cfg)).Directory(templateDir, ".html")

	web.Get("/", index)
	web.Get("/index.html", index)
	web.Get("/search", search)

	web.StaticWeb("/static", g.Conf.StaticDirectory, 1)

	web.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		ctx.Render("500.html", nil)
	})
	web.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Render("404.html", nil)
	})

	web.Listen(g.Conf.Address)
}
