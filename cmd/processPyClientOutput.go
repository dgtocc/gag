package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func ProcessPyClientOutput(f string) error {

	var tstypemapper map[string]string = make(map[string]string)
	var exceptionaltypemapper map[string]string = make(map[string]string)
	exceptionaltypemapper["[]byte"] = "str"
	exceptionaltypemapper["[]string"] = "List[str]"

	tstypemapper["string"] = "str"
	tstypemapper["time.Time"] = "int"
	tstypemapper["primitive.ObjectID"] = "str"
	tstypemapper["time.Duration"] = "int"
	tstypemapper["int"] = "int"
	tstypemapper["int32"] = "int"
	tstypemapper["int64"] = "int"
	tstypemapper["float"] = "float"
	tstypemapper["float64"] = "float"
	tstypemapper["uint8"] = "int"
	tstypemapper["uint16"] = "int"
	tstypemapper["uint32"] = "int"
	tstypemapper["error"] = "Exception"
	tstypemapper["bool"] = "bool"
	tstypemapper["interface{}"] = "dict"
	tstypemapper["bson.M"] = "dict"

	_typeName := func(m *APIParamType) string {
		if m.IsArray {
			return "List[" + m.Typename + "]"
		}
		return m.Typename
	}

	b := bytes.Buffer{}

	b.WriteString("from dataclasses import dataclass\n")
	b.WriteString("import requests\n")
	b.WriteString("import json\n")
	b.WriteString("from typing import List\n")

	b.WriteString("#region Base\n")
	b.WriteString(fmt.Sprintf(`

__ctx = {"apibase":"%s"}


def SetAPIBase(s: str):
    __ctx["apibase"] = s


def GetAPIBase() -> str:
    return __ctx["apibase"]


def SetCookie(s: str):
    __ctx["cookie"] = s


def GetCookie() -> str:
    return __ctx["cookie"]


def InvokeTxt(path: str, method: str, body) -> str:
    headers = {"Content-type": "application/json", "Cookie": "dc="+GetCookie()}
    fpath = GetAPIBase() + path
    r = requests.request(method, fpath, json=body, headers=headers)
    return r.text


def InvokeJSON(path: str, method: str, body) -> dict:
	d = body.__dict__
	return json.loads(InvokeTxt(path, method, d))

`, api.BasePath))
	b.WriteString("#endregion\n\n")
	b.WriteString("#region Types\n")
	for k, v := range api.Types {
		if v.Desc != "" {
			//b.WriteString(fmt.Sprintf("/**\n%s*/\n", v.Desc))
		}
		if len(v.Fields) < 1 {
			b.WriteString(fmt.Sprintf("@ dataclass\nclass %s :\n\tpass\n\n", k))
		} else {
			b.WriteString(fmt.Sprintf("@ dataclass\nclass %s :\n", k))
			var ftype string
			var ok bool
			for kf, f := range v.Fields {
				ftype, ok = exceptionaltypemapper[f.String()]
				if ok {
					log.Printf("Mapped exceptional type: %s ==> %s", f.String(), ftype)
				}

				if !ok {
					if f.Array {
						ftype, ok = tstypemapper["[]"+f.Type]
					} else {
						ftype, ok = tstypemapper[f.Type]
					}
				}

				if !ok {
					ftype = f.Type
				}
				if f.Map {
					//fm, ok := tstypemapper[f.Mapkey]
					//if !ok {
					//	fm = f.Mapkey
					//}
					//fv, ok := tstypemapper[f.Mapval]
					//if !ok {
					//	fv = f.Mapval
					//}
					//ftype = "{[s:" + fm + "]:" + fv + "}"
					ftype = "dict"
				}

				if f.Desc != "" {
					//b.WriteString(fmt.Sprintf("\t/**\n%s*/\n", f.Desc))
				}
				b.WriteString(fmt.Sprintf("\t%s: %s\n", strings.ToLower(kf), ftype))
			}

			b.WriteString(fmt.Sprintf("\n\n"))
		}
	}
	b.WriteString("#endregion\n\n")
	b.WriteString("#region Methods\n")
	for k, m := range api.Methods {
		if m.Desc != "" {
			//b.WriteString(fmt.Sprintf("/**\n%s*/\n", m.Desc))
		}

		rettype := _typeName(m.ResType)
		if rettype != "" {
			b.WriteString(fmt.Sprintf("def %s(req:%s)-> %s:\n", k, _typeName(m.ReqType), rettype))
		} else {
			b.WriteString(fmt.Sprintf("def %s(req:%s):\n", k, _typeName(m.ReqType)))
		}

		b.WriteString(fmt.Sprintf("\treturn InvokeJSON(\"%s\",\"%s\",req)\n", m.Path, m.Verb))
		b.WriteString(fmt.Sprintf("\n\n"))
		//}

	}
	b.WriteString("#endregion\n")

	err := ioutil.WriteFile(f, b.Bytes(), 0600)
	return err
}
