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
