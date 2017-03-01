package controllers

import (
    //"github.com/QLeelulu/goku"
	"goku-master/goku-master"
)

var _ = goku.Controller("home").
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
		return ctx.Html("html1");
    //return ctx.Redirect("/")
})
