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

Based on the metadata added to your method, the API will be generated

## Install

```
go install github.com/dgtocc/gag@latest
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
