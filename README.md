# GAG
Golang API Generator

Hi, hate bureaucracy? Long and verbose approach for simple things? So do I... 

In this case, if you want to derive API Servers and Clients from your functions in GO, and as me, you are tired of 
endless YAMLS or JSONS, please have a sit and join the gang.

GAG is all about reducing the friction for making your GO code available and consumabe as API.

## How it works
Suppose you want to generate an API Server, using gin as backend:

create your GO package, and add your API Methods like this:

```go
/*
@API
@PATH: /someapi
@PERM: ASD
@VERB: POST
*/
func SomeAPI(ctx context.Context, s string) (out string, err error) {
	print("Got:" + s)
	out = time.Now().String() + " - Hey Ya!"
	return
}
```

> Important: All Methods must have the same pattern: 2 params, 1st being the context, 2nd the request object. 
> Should return 2 params also. 1st the response object, 2nd an eventual error 


Suppose this is saved as: `/home/me/project/src/api/some.go`

Now, run gag like: 

`gag gin /home/me/project/src/api`

You will find a new file generated there, called apigen with all the dirty work in it.

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
User defined structs can be exported also:

```go

/*@API*/
type AStr struct {
	Country       string
	City          string
	HouseNumber   int64
	IsCondo       bool
	SomeWeirdTest string `json:"SUPERCALIFRAGILISPEALIDOUX"`
	Recursive     map[string]AStr
	Arrofpstr     []string `json:"arrofpstr,omitempty"`
	When          time.Time
	Some          crypto.Decrypter
}


/*
@API
@PATH: /someapi2
@PERM: ASD
@VERB: GET
*/
func SomeGET(ctx context.Context, s []*AStr) (out string, err error) {
	print("Got:" + s[0].SomeWeirdTest)
	out = time.Now().String() + " - Hey Ya!"
	return
}
```


Based on the metadata added to your method, the API will be generated

In case you have not noticed, comments are metadata for the API:

- API is a marker, both methods and structs are exported only in case this 
tag is found.
- PATH: applies to methods, and in case they exist, bind the method to path. In case absent,method name becomes path, 
  having _ to be replaced by / (func a_b_c will be bound to path a/b/c)
- PERM: defines permission group to method - please see later how to bind permissions
- VERB: Http Verb to be used when binding method. Default is POST

## Install

```
go build -o gag ./cmd
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
