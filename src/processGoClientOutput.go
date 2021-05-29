package main

import (
	"fmt"
	"os"
	"strings"
)

import (
	_ "embed"
)

func processGoClientOutput(f string) error {
	sb := strings.Builder{}
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

	l := func(s string, p ...interface{}) {
		sb.WriteString(fmt.Sprintf(s+"\n", p...))
	}
	l("package %s", api.Namespace)
	l(`import (
	"bytes"
	"errors"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"
)
var Basepath string = ""
var Host string = ""
var ExtraHeaders map[string]string = make(map[string]string)

func invoke(m string, path string, bodyo interface{}) (*json.Decoder, error) {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(bodyo)
	if err != nil {
		return nil, err
	}
		
	body := bytes.NewReader(b.Bytes())
	req, err := http.NewRequest(m, Host+Basepath+path, body)
		
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")

	for k, v := range ExtraHeaders {
		req.Header.Set(k, v)
	}

	cli := http.Client{}
	res, err := cli.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		return nil, errors.New(string(bs))
	}

	ret := json.NewDecoder(res.Body)
	return ret, nil
}
`)
	for tname, t := range api.Types {
		l("type %s struct{", tname)
		for fname, f := range t.Fields {
			if f.Map {
				l("\t %s map[%s]%s", fname, f.Mapkey, f.Mapval)
			} else if f.Array {
				l("\t %s[]", f.Type)
			} else {
				l("\t %s", f.Type)
			}
		}
		l("}", tname)
	}

	for _, m := range api.Methods {
		l("func %s(req %s)(res %s,err error){", m.Name, resolveReqTypeDecStr(m.ReqType), resolveReqTypeDecStr(m.ResType))
		l("\tvar dec *json.Decoder")
		l("\tdec,err=invoke(\"%s\",\"%s\",req)", m.Verb, m.Path)
		l("\tif err!=nil{")
		l("\t\treturn")
		l("\t}")
		l("\tvar ret %s", resolveReqTypeDecStr(m.ResType))
		amp := ""
		if m.ResType.IsArray || !m.ResType.Ispointer {
			amp = "&"
		}
		l("\terr=dec.Decode(%sret)", amp)
		l("return ret,err")
		l("}")
	}

	return os.WriteFile(f, []byte(sb.String()), 0600)
}
