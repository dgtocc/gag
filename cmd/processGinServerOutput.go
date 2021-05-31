package main

import (
	"fmt"
	"os"
	"strings"
)

func ProcessGinServerOutput(f string) error {
	sb := strings.Builder{}
	l := func(s string, p ...interface{}) {
		if !strings.HasSuffix(s, "\n") {
			s = s + "\n"
		}
		sb.WriteString(fmt.Sprintf(s, p...))
	}

	resolveReqTypeDecStr := func(r *APIParamType) string {
		ret := ""
		if r.IsArray {
			ret = "[]"
		}
		if r.Ispointer {
			ret = ret + "*"
		}
		ret = ret + r.Typename
		return ret
	}

	resolveReqTypeInstStr := func(r *APIParamType) string {
		sb := strings.Builder{}

		if r.IsArray {
			sb.WriteString("make([]")
			if r.Ispointer {
				sb.WriteString("*")
			}
			sb.WriteString(r.Typename)
			sb.WriteString(",0)")
		} else {
			if r.Ispointer {
				sb.WriteString("&")
			}
			sb.WriteString(r.Typename)
		}

		return sb.String()
	}

	//tmpl, err := template.New("gin").Parse(ginServerTemplate)

	l(fmt.Sprintf("package %s\n\n", api.Namespace))
	l("import(")
	l("\t\"github.com/gin-gonic/gin\"")
	l("\tcontextlib \"context\"")
	l("\t\"net/http\"")

	for alias, pkg := range api.UsedImportsFunctions {
		if alias != "net/http" {
			if alias != pkg {
				l(fmt.Sprintf("\t%s \"%s\"\n", alias, pkg))
			} else {
				l(fmt.Sprintf("\t\"%s\"\n", pkg))
			}
		}
	}

	l(")")
	l("var perms map[string]string")
	l("func init(){")
	l("\tperms=make(map[string]string)")
	for _, m := range api.Methods {
		if m.Perm != "" {
			l(fmt.Sprintf("	perms[\"%s_%s\"]=\"%s\"\n", m.Verb, m.Path, m.Perm))
		}
	}
	l("}")
	l(`func GetPerm(c *gin.Context) string {
	perm, ok := perms[c.Request.Method+"_"+c.Request.URL.Path]
	if !ok {
		return ""
	}
	return perm
}

`)
	l("func Build(r *gin.Engine) {")
	for _, path := range api.SortedPaths {
		for _, verb := range path.SortedVerbs {
			l("\tr.%s(\"%s\",func(c *gin.Context) {", verb.Verb, path.Path)
			l("\t\tvar req %s", resolveReqTypeDecStr(verb.Method.ReqType))
			if verb.Method.ReqType.IsArray {
				if verb.Method.ReqType.Ispointer {
					l("\t\treq = make([]*%s,0)", verb.Method.ReqType.Typename)
				} else {
					l("\t\treq = make([]%s,0)", verb.Method.ReqType.Typename)
				}
				l("\t\tc.BindJSON(&req)")
			} else if verb.Method.ReqType.Ispointer {
				l("\t\treq = %s{}", resolveReqTypeInstStr(verb.Method.ReqType))
				l("\t\tc.BindJSON(req)")
			} else {
				l("\t\tc.BindJSON(&req)")
			}
			l("\t\tc.Request.WithContext(contextlib.WithValue(c.Request.Context(), \"CTX\", c))")
			l("\t\tc.Request.WithContext(contextlib.WithValue(c.Request.Context(), \"REQ\", c.Request))")
			l("\t\tc.Request.WithContext(contextlib.WithValue(c.Request.Context(), \"RES\", c.Writer))")
			l("\t\tres,err := %s(c.Request.Context(),req)", verb.Method.Name)
			l("\t\tif err!=nil{")
			l("\t\t\tc.AbortWithError(http.StatusInternalServerError,err)")
			l("\t\t\treturn")
			l("\t\t}")
			l("\t\tc.JSON(200,res)")
			l("\t})")
		}
	}
	l("}")

	return os.WriteFile(f, []byte(sb.String()), 0600)
}
