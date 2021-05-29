package main

import (
	_ "embed"
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
		sb.WriteString(fmt.Sprintf(s+"\n", p...))
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
	l("import(\n")
	l("\t\"github.com/gin-gonic/gin\"\n")
	l("\t\"net/http\"\n")

	for alias, pkg := range api.UsedImportsFunctions {
		if alias != "net/http" {
			if alias != pkg {
				l(fmt.Sprintf("\t%s \"%s\"\n", alias, pkg))
			} else {
				l(fmt.Sprintf("\t\"%s\"\n", pkg))
			}
		}
	}

	l(")\n")
	l("var perms map[string]string\n\n")
	l("func init(){\n")
	l("\tperms=make(map[string]string)\n")
	for _, m := range api.Methods {
		if m.Perm != "" {
			l(fmt.Sprintf("	perms[\"%s_%s\"]=\"%s\"\n", m.Verb, m.Path, m.Perm))
		}
	}
	l("}\n\n")
	l(`func GetPerm(c *gin.Context) string {
	perm, ok := perms[c.Request.Method+"_"+c.Request.URL.Path]
	if !ok {
		return ""
	}
	return perm
}

`)
	l("func Build(r *gin.Engine) {\n")
	for _, path := range api.SortedPaths {
		for _, verb := range path.SortedVerbs {
			l(fmt.Sprintf("\tr.%s(\"%s\",func(c *gin.Context) {\n", verb.Verb, path.Path))
			l(fmt.Sprintf("\t\tvar req %s\n", resolveReqTypeDecStr(verb.Method.ReqType)))
			if verb.Method.ReqType.Ispointer {
				l(fmt.Sprintf("\t\treq = %s\n", resolveReqTypeInstStr(verb.Method.ReqType)))
				l("\t\tc.BindJSON(req)\n")
			} else {
				l("\t\tc.BindJSON(&req)\n")
			}
			l(fmt.Sprintf("\t\tres,err := %s(c.Request.Context(),req)\n", verb.Method.Name))
			l("\t\tif err!=nil{\n")
			l("\t\t\tc.AbortWithError(http.StatusInternalServerError,err)\n")
			l("\t\t\treturn\n")
			l("\t\t}\n")
			l("\t\tc.JSON(200,res)\n")
			l("\t})\n\n")
		}
	}
	l("}\n")

	return os.WriteFile(f, []byte(sb.String()), 0600)
}
