package goapi

import (
	contextlib "context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var perms map[string]string

func init() {
	perms = make(map[string]string)
	perms["GET_/someapi2"] = "ASD"
	perms["POST_/someapi"] = "ASD"
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
		var req []time.Time
		req = make([]time.Time, 0)
		c.BindJSON(&req)
		c.Request.WithContext(contextlib.WithValue(c.Request.Context(), "CTX", c))
		res, err := SomeAPI(c.Request.Context(), req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, res)
	})
	r.GET("/someapi2", func(c *gin.Context) {
		var req []*AStr
		req = make([]*AStr, 0)
		c.BindJSON(&req)
		res, err := SomeGET(c.Request.Context(), req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, res)
	})
}
