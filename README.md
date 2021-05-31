# GAG

Golang API Generator

Hi, hate bureaucracy? Long and verbose approach for simple things? So do I...

In this case, if you want to derive API Servers and Clients from your functions in GO, and as me, you are tired of
endless YAMLS or JSONS, please have a sit and join the gang.

GAG is all about reducing the friction for making your GO code available and consumabe as API.

## How it works

Suppose you want to generate an API Server, using gin as backend:

Create you business logic as shown below:

```go
package goapi

import (
	"context"
	"crypto"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

/*@API*/
type ASimpleReq struct {
	Data string
}

/*@API*/
type ASimpleRes struct {
	Data string
}

/*@API*/
type AComplexReq struct {
	Country       string
	City          string
	HouseNumber   int64
	IsCondo       bool
	SomeWeirdTest string `json:"SUPERCALIFRAGILISPEALIDOUX"`
	Recursive     map[string]AComplexReq
	Arrofpstr     []string `json:"arrofpstr,omitempty"`
	When          time.Time
	Some          crypto.Decrypter
}

/*
@API
@PATH: /someapi
@PERM: ASD
@VERB: POST
*/
func ApiMethod01(ctx context.Context, s *ASimpleReq) (out *ASimpleRes, err error) {
	log.Printf("Got: %#v", s)
	out = &ASimpleRes{}
	out.Data = time.Now().String() + " - Hey Ya!"
	return
}

/*
@API
@PATH: /someapi2
@PERM: ASD
@VERB: GET
*/
func ApiMethod02(ctx context.Context, s *AComplexReq) (out *ASimpleRes, err error) {
	gctx := ctx.Value("CTX").(*gin.Context)
	log.Printf(gctx.FullPath())
	print("Got:" + s.SomeWeirdTest)
	out = &ASimpleRes{}
	out.Data = time.Now().String() + " - Hey Ya!"
	return
}

```

> Important: All Methods must have the same pattern: 2 params, 1st being the context, 2nd the request object.
> Should return 2 params also. 1st the response object, 2nd an eventual error.
>
> **ALTHOUGHT you can use virtually any type as request and response we propose that you always use a POINTER TO a data structure
> as request and response - dont be creative, be practical**




Suppose this is saved as: `/home/me/project/src/api/some.go`

Now, run gag like:

`gag gin /home/me/project/src/api`

You will find a new file generated there, called apigen with all the dirty work in it.

```go
package goapi

import (
	contextlib "context"
	"github.com/gin-gonic/gin"
	"net/http"
)

var perms map[string]string

func init() {
	perms = make(map[string]string)
	perms["POST_/someapi"] = "ASD"
	perms["GET_/someapi2"] = "ASD"
}
func GetPerm(c *gin.Context) string {
	perm, ok := perms[c.Request.Method+"_"+c.Request.URL.Path]
	if !ok {
		return ""
	}
	return perm
}

func Build(r *gin.Engine) {
	r.POST("/someapi", func(c *gin.Context) {
		var req *ASimpleReq
		req = &ASimpleReq{}
		c.BindJSON(req)
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "CTX", c))
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "REQ", c.Request))
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "RES", c.Writer))
		res, err := ApiMethod01(c.Request.Context(), req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, res)
	})
	r.GET("/someapi2", func(c *gin.Context) {
		var req *AComplexReq
		req = &AComplexReq{}
		c.BindJSON(req)
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "CTX", c))
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "REQ", c.Request))
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "RES", c.Writer))
		res, err := ApiMethod02(c.Request.Context(), req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, res)
	})
}


```

Later, you may create your server main like: `/home/me/project/src/api/some.go` and make it like:

```go
package main

import (
	"github.com/dgtocc/gag/test/goapi"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	goapi.Build(g)
	g.Run()
}

```

Please pay attention to data structures - they also need to be tagged as API to be exported

## HOW IT WORKS

Based on the metadata added to your method, the API will be generated

In case you have not noticed, comments are metadata for the API:

- API is a marker, both methods and structs are exported only in case this tag is found.
- PATH: applies to methods, and in case they exist, bind the method to path. In case absent,method name becomes path,
  having _ to be replaced by / (func a_b_c will be bound to path a/b/c)
- PERM: defines permission group to method - please see later how to bind permissions
- VERB: Http Verb to be used when binding method. Default is POST

## Pluging in the permissions

Please see the method checkpermisson below. Once you generate the server, the method GetPerm will be created, and can be
used to retrieve the permission set for a call. By usint a middleware like this you can easily plug any permission
system available with the generated API.

```go
package main

import (
	"github.com/dgtocc/gag/test/goapi"
	"github.com/gin-gonic/gin"
	"net/http"
)

func checkPerm(ctx *gin.Context, perm string) bool {
	//Do your logic here
	return false
}
func main() {
	g := gin.Default()
	//Setup middlewares before processing API
	g.Use(func(context *gin.Context) {
		perm := goapi.GetPerm(context)
		if checkPerm(context, perm) == false {
			context.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			context.Next()
		}
	})

	//Bind built API
	goapi.Build(g)

	//Run, baby run...
	g.Run()
}

```

## Oh, but I still need to access something from GIN or HTTP Context,Request, Response....

Worry not - context follows your code. At any time you can call it like:

```go
gctx:= ctx.Value("CTX").(*gin.Context)
req:= ctx.Value("REQ").(http.Request)
res:= ctx.Value("RES").(http.ResponseWriter)
```

Keep in mind that, if new server implementation is created such dependency will break if you try to generate such code.
You can use middlewares for this in your code, if it is the case.

## Install

Grab the binaries :D or,

```
make release
```

## General Usage

```
Usage: gag <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  yaml <src> <fname>
    Gens YAML metamodel

  gin <src>
    Gens Gin Server impl

  gocli <src> <dst>
    Gens Go Cli impl

  pycli <src> <dst>
    Gens Python Cli impl

  ts <src> <dst>
    Gens Typescript Cli impl

  http <src> <dst>
    Gens Http call impl

```

## Need improvement

- Type mapping system: Here the doubt is more on the practical side. It seems to be good enought
  for the basic request and response data objects. Need inputs on concrete requirements for this improvement
- loader code is too messy: Lets see...
- No integration with lower level HTTP protocol: Ok - this is a high level abstraction tool. This can be overcome by either
using objects from context or through middlewares.
  
Feedback on these is highly appreciated.
