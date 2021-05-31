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
