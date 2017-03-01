#goku

  goku is a Web Mvc Framework for golang, mostly like ASP.NET MVC.

  [doc & api](http://go.pkgdoc.org/github.com/QLeelulu/goku)

##Installation

To install goku, simply run `go get github.com/QLeelulu/goku`. To use it in a program, use `import "github.com/QLeelulu/goku"`

To run example "todo" app, just:
    
    $ cd $GOROOT/src/pkg/github.com/QLeelulu/goku/examples/todo/
    $ go run app.go

maybe you need run todo.sql first.

##Usage

```go
    package main

    import (
        "github.com/QLeelulu/goku"
        "log"
        "path"
        "runtime"
        "time"
    )

    // routes
    var routes []*goku.Route = []*goku.Route{
        // static file route
        &goku.Route{
            Name:     "static",
            IsStatic: true,
            Pattern:  "/static/(.*)",
        },
        // default controller and action route
        &goku.Route{
            Name:       "default",
            Pattern:    "/{controller}/{action}/{id}",
            Default:    map[string]string{"controller": "home", "action": "index", "id": "0"},
            Constraint: map[string]string{"id": "\\d+"},
        },
    }

    // server config
    var config *goku.ServerConfig = &goku.ServerConfig{
        Addr:           ":8888",
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
        //RootDir:        os.Getwd(),
        StaticPath: "static",
        ViewPath:   "views",
        Debug:      true,
    }

    func init() {
        /**
         * project root dir
         */
        _, filename, _, _ := runtime.Caller(1)
        config.RootDir = path.Dir(filename)

        /**
         * Controller & Action
         */
        goku.Controller("home").
            Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
            return ctx.Html("Hello World")
        })

    }

    func main() {
        rt := &goku.RouteTable{Routes: routes}
        s := goku.CreateServer(rt, nil, config)
        goku.Logger().Logln("Server start on", s.Addr)
        log.Fatal(s.ListenAndServe())
    }
```


##Route
    
```go
    var Routes []*goku.Route = []*goku.Route{
        &goku.Route{
            Name:     "static",
            IsStatic: true,
            Pattern:  "/public/(.*)",
        },
        &goku.Route{
            Name:       "edit",
            Pattern:    "/{controller}/{id}/{action}",
            Default:    map[string]string{"action": "edit"},
            Constraint: map[string]string{"id": "\\d+"},
        },
        &goku.Route{
            Name:    "default",
            Pattern: "/{controller}/{action}",
            Default: map[string]string{"controller": "todo", "action": "index"},
        },
    }
    // create server with the rules
    rt := &goku.RouteTable{Routes: todo.Routes}
    s := goku.CreateServer(rt, nil, nil)
    log.Fatal(s.ListenAndServe())
```

or

```go
    rt := new(goku.RouteTable)
    rt.Static("staticFile", "/static/(.*)")
    rt.Map(
        "blog-page", // route name
        "/blog/p/{page}", // pattern
        map[string]string{"controller": "blog", "action": "page", "page": "0"}, // default
        map[string]string{"page": "\\d+"} // constraint
    )
    rt.Map(
        "default", // route name
        "/{controller}/{action}", // pattern
        map[string]string{"controller": "home", "action": "index"}, // default
    )
```

##Controller And Action

```go
    // add a controller named "home"
    goku.Controller("home").
        Filters(new(TestControllerFilter)). // this filter is fot controller(all the actions)
        // home controller's index action
        // for http GET
        Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.Html("Hello World")
    }).Filters(new(TestActionFilter)). // this filter is for home.index action
        // home controller's about action
        // for http POST
        Post("about", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.Raw("About")
    })
```

+ get  `/home/index` will return `Hello World`
+ post `/home/index` will return 404
+ post `/home/about` will return `About`

##ActionResult

`ActionResulter` is type interface. all the action must return ActionResulter.
you can return an ActionResulter in the action like this:

```go
    // ctx is *goku.HttpContext
    ctx.Raw("hi")
    ctx.NotFound("oh no! ):")
    ctx.Redirect("/")
    // or you can return a view that
    // will render a template
    ctx.View(viewModel)
    ctx.Render("viewName", viewModel)
```

or you can return a ActionResulter by this

```go
    return &ActionResult{
        StatusCode: http.StatusNotFound,
        Headers:    map[string]string{"Content-Type": "text/html"},
        Body:       "Page Not Found",
    }
```

for more info, check the code.

##View and ViewData

`Views are the components that display the application's user interface (UI)`

To render a view, you can just return a `ViewResut` which implement the `ActionResulter` interface.
just like this:

```go
    goku.Controller("blog").
        Get("page", func(ctx *goku.HttpContext) goku.ActionResulter {
        blog := GetBlogByid(10)

        // you can add any val to ViewData
        // then you can use it in template
        // like this: {{ .Data.SiteName }}
        ctx.ViewData["SiteName"] = "My Blog"

        // or you can pass a struct to ViewModel
        // then you can use it in template
        // like this: {{ .Model.Title }}
        // that same as blog.Title
        return ctx.View(blog)
    })
```

`ctx.View()` will find the view in these rules:
  
1. /{ViewPath}/{Controller}/{action}
2. /{ViewPath}/shared/{action}

for example, ServerConfig.ViewPath is set to "views",
and return ctx.View() in `home` controller's `about` action, 
it will find the view file in this rule:

1. {ProjectDir}/views/home/about.html
2. {ProjectDir}/views/shared/about.html

if you want to return a view that specified view name,
you can use ctx.Render:

```go
    // it will find the view in these rules:
    //      1. /{ViewPath}/{Controller}/{viewName}
    //      2. /{ViewPath}/shared/{viewName}
    // if viewName start with '/',
    // it will find the view direct by viewpath:
    //      1. /{ViewPath}/{viewName}
    ctx.Render("viewName", ViewModel)
```

##ViewEngine & Template

```go
    // you can add any val to ViewData
    // then you can use it in template
    // like this: {{ .Data.SiteName }}
    ctx.ViewData["SiteName"] = "My Blog"

    blogs := GetBlogs()
    // or you can pass a struct to ViewModel
    // then you can use it in template
    // like this: {{range .Model}} {{ .Title }} {{end}}
    return ctx.View(blogs)
```

default template engine is golang's template.

```html
    <div class="box todos">
        <h2 class="box">{{ .Data.SiteName }}</h2>
        <ul>
          {{range .Model}}
            <li id="blog-{{.Id}}">
              {{.Title}}
            </li>
          {{end}}
        </ul>
    </div>
```

####Layout

layout.html
```html
    <!DOCTYPE html>
    <html>
    <head>
        <title>Goku</title>
        {{template "head"}}
    </head>
    <body>
      {{template "body" .}}
    </body>
    </html>
```

body.html
```html
    {{define "head"}}
        <!-- add css or js here -->
    {{end}}

    {{define "body"}}
        I'm main content.
    {{end}}
```
note the dot in `{{template "body" .}}` , it will pass the ViewData to the sub template.

HtmlHelper?

####More Template Engine Support

if you want to use [mustache](https://github.com/hoisie/mustache) template, 
check [mustache.goku](https://github.com/QLeelulu/mustache.goku)


##HttpContext

```go
    type HttpContext struct {
        Request        *http.Request       // http request
        responseWriter http.ResponseWriter // http response
        Method         string              // http method

        //self fields
        RouteData *RouteData             // route data
        ViewData  map[string]interface{} // view data for template
        Data      map[string]interface{} // data for httpcontex
        Result    ActionResulter         // action result
        Err       error                  // process error
        User      string                 // user name
        Canceled  bool                   // cancel continue process the request and return
    }
```


#Form Validation

you can create a form, to valid the user's input, and get the clean value.

```go
    import "github.com/QLeelulu/goku/form"

    func CreateCommentForm() *goku.Form {
        name := NewCharField("name", "Name", true).Range(3, 10).Field()
        nickName := NewCharField("nick_name", "Nick Name", false).Min(3).Max(20).Field()
        age := NewIntegerField("age", "Age", true).Range(18, 50).Field()
        content := NewTextField("content", "Content", true).Min(10).Field()

        form := NewForm(name, nickName, age, content)
        return form
    }
```

and then you can use this form like this:

```go
    f := CreateCommentForm()
    f.FillByRequest(ctx.Request)

    if f.Valid() {
        // after valid, we can get the clean values
        m := f.CleanValues()
        // and now you can save m to database
    } else {
        // if not valid
        // we can get the valid errors
        errs := f.Errors()
    }
```

checkout [form_test.go](https://github.com/QLeelulu/goku/blob/master/form/form_test.go)


##DataBase

simple database api.

```go
    db, err := OpenMysql("mymysql", "tcp:localhost:3306*test_db/lulu/123456")

    // you can use all the api in golang's database/sql
    _, err = db.Query("select 1")

    // or you can use some simple api provide by goku
    r, err := db.Select("test_blog", SqlQueryInfo{
        Fields: "id, title, content",
        Where:  "id>?",
        Params: []interface{}{0},
        Limit:  10,
        Offset: 0,
        Group:  "",
        Order:  "id desc",
    })

    // insert map
    vals := map[string]interface{}{
        "title": "golang",
        "content": "Go is an open source programming environment that " +
            "makes it easy to build simple, reliable, and efficient software.",
        "create_at": time.Now(),
    }
    r, err := db.Insert("test_blog", vals)

    // insert struct
    blog := TestBlog{
        Title:    "goku",
        Content:  "a mvc framework",
        CreateAt: time.Now(),
    }
    r, err = db.InsertStruct(&blog)

    // get struct
    blog := &TestBlog{}
    err = db.GetStruct(blog, "id=?", 3)

    // get struct list
    qi := SqlQueryInfo{}
    var blogs []Blog
    err := db.GetStructs(&blogs, qi)

    // update by map
    vals := map[string]interface{}{
        "title": "js",
    }
    r, err2 := db.Update("test_blog", vals, "id=?", blog.Id)

    // delete
    r, err := db.Delete("test_blog", "id=?", 8)
```

checkout [db_test.go](https://github.com/QLeelulu/goku/blob/master/db_test.go)

####DataBase SQL Debug

if you want to debug what the sql query is, set db.Debug to `true`

```go
    db, err := OpenMysql("mymysql", "tcp:localhost:3306*test_db/username/pwd")
    db.Debug = true
```

after you set db.Debug to true, while you run a db command, it will print the sql query to the log,
juse like this:

    2012/07/30 20:58:03 SQL: UPDATE `user` SET friends=friends+? WHERE id=?;
                        PARAMS: [[1 2]]


##Action Filter

```go
    type TestActionFilter struct {
    }

    func (tf *TestActionFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
        ctx.WriteString("OnActionExecuting - TestActionFilter \n")
        return
    }
    func (tf *TestActionFilter) OnActionExecuted(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
        ctx.WriteString("OnActionExecuted - TestActionFilter \n")
        return
    }

    func (tf *TestActionFilter) OnResultExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
        ctx.WriteString("    OnResultExecuting - TestActionFilter \n")
        return
    }

    func (tf *TestActionFilter) OnResultExecuted(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
        ctx.WriteString("    OnResultExecuted - TestActionFilter \n")
        return
    }
```

Order of the filters execution is:

   1. OnActionExecuting
   2. -> Execute Action -> return ActionResulter
   3. OnActionExecuted
   4. OnResultExecuting
   5. -> ActionResulter.ExecuteResult
   6. OnResultExecuted



##Middleware

```go
    type TestMiddleware struct {
    }

    func (tmd *TestMiddleware) OnBeginRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
        ctx.WriteString("OnBeginRequest - TestMiddleware \n")
        return nil, nil
    }

    func (tmd *TestMiddleware) OnBeginMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
        ctx.WriteString("  OnBeginMvcHandle - TestMiddleware \n")
        return nil, nil
    }
    func (tmd *TestMiddleware) OnEndMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
        ctx.WriteString("  OnEndMvcHandle - TestMiddleware \n")
        return nil, nil
    }

    func (tmd *TestMiddleware) OnEndRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
        ctx.WriteString("OnEndRequest - TestMiddleware \n")
        return nil, nil
    }
```

Order of the middleware event execution is:

   1. `OnBeginRequest`
   2. `OnBeginMvcHandle`(if not the static file request)
   3.  => run controller action (if not the static file request)
   4. `OnEndMvcHandle`(if not the static file request)
   5. `OnEndRequest`


##Log

To use logger in goku, just:

    goku.Logger().Logln("i", "am", "log")
    goku.Logger().Errorln("oh", "no!", "Server Down!")

this will log like this:

    2012/07/14 20:07:46 i am log
    2012/07/14 20:07:46 [ERROR] oh no! Server Down!

## Authors

 - [QLeelulu](https://github.com/QLeelulu)
 - waiting for you


## License

View the [LICENSE](https://github.com/QLeelulu/goku/blob/master/LICENSE) file.