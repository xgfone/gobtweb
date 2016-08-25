package main

import (
	"os"
	"strings"

	"github.com/iris-contrib/template/django"
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/xgfone/gobtweb/g"
)

func handlerError(ctx *iris.Context, err error) {
	if err != nil {
		g.Logger.Error("Failed to return the response", "err", err)
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
		ctx.Redirect("/")
		return
	}

	var from int
	var size int
	var err error

	if from, err = ctx.URLParamInt("from"); err != nil {
		from = 0
	}

	if size, err = ctx.URLParamInt("size"); err != nil {
		size = 10
	}

	queries := make([]string, 0)
	for _, q := range strings.Split(query, " ") {
		q = strings.TrimSpace(q)
		if q != "" {
			queries = append(queries, q)
		}
	}

	data, err := SearchKeyword(g.ElasticClient, queries, from, size)
	if err != nil {
		handlerError(ctx, err)
		return
	}

	err = ctx.Render("search.html", map[string]interface{}{
		"query":    query,
		"torrents": data,
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
	web.UseTemplate(django.New(django_cfg)).Directory(templateDir, ".html")

	web.Get("/", index)
	web.Get("/index.html", index)
	web.Get("/search", search)

	web.StaticWeb("/static", g.Conf.StaticDirectory, 1)

	web.Listen(g.Conf.Address)
}
