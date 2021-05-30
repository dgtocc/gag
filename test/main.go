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
	g.Use(func(context *gin.Context) {
		perm := goapi.GetPerm(context)
		if checkPerm(context, perm) == false {
			context.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			context.Next()
		}
	})
	g.Run()
}
