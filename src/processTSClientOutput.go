package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func processTSClientOutput(f string) error {

	var tstypemapper map[string]string = make(map[string]string)
	var exceptionaltypemapper map[string]string = make(map[string]string)
	exceptionaltypemapper["[]byte"] = "string"

	tstypemapper["time.Time"] = "Date"
	tstypemapper["primitive.ObjectID"] = "string"
	tstypemapper["time.Duration"] = "Date"
	tstypemapper["int"] = "number"
	tstypemapper["int32"] = "number"
	tstypemapper["int64"] = "number"
	tstypemapper["float"] = "number"
	tstypemapper["float64"] = "number"
	tstypemapper["uint8"] = "number"
	tstypemapper["uint16"] = "number"
	tstypemapper["uint32"] = "number"
	tstypemapper["error"] = "Error"
	tstypemapper["bool"] = "boolean"
	tstypemapper["interface{}"] = "any"
	tstypemapper["bson.M"] = "any"

	_typeName := func(m *APIParamType) string {
		if m.IsArray {
			return m.Typename + "[]"
		}
		return m.Typename
	}

	b := bytes.Buffer{}

	b.WriteString("//#region Base\n")
	b.WriteString(fmt.Sprintf(`

var apibase="%s";

export function SetAPIBase(s:string){
	apibase=s;
}

export function GetAPIBase(): string{
	return apibase;
}

let REGEX_DATE = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2}(?:\.\d*)?)(Z|([+\-])(\d{2}):(\d{2}))$/

type HTMLMethod = "GET" | "POST" | "PUT" | "DELETE" | "HEAD" | "TRACE"

async function Invoke(path: string, method: HTMLMethod, body?: any): Promise<Response> {
	let jbody = undefined
	let init = {method: method, mode: "cors", credentials: "include", withCredentials: true}
	if (!!body) {
		let jbody = JSON.stringify(body)
		//@ts-ignore
		init.body = jbody
	}
	if (apibase.endsWith("/") && path.startsWith("/")) {
		path = path.substr(1, path.length)
	}
	let fpath = (apibase + path)
	//@ts-ignore
	let res = await fetch(fpath, init)

	return res
}
 
async function InvokeJSON(path: string, method: HTMLMethod, body?: any): Promise<any> {

	let txt = await InvokeTxt(path, method, body)
	if (txt == "") {
		txt = "{}"
	}
	let ret = JSON.parse(txt, (k: string, v: string) => {
		if (REGEX_DATE.exec(v)) {
			return new Date(v)
		}
		return v
	})

	return ret
}

async function InvokeTxt(path: string, method: HTMLMethod, body?: any): Promise<string> {
	//@ts-ignore
	let res = await Invoke(path, method, body)

	let txt = await res.text()

	if (res.status < 200 || res.status >= 400) {
		// webix.alert("API Error:" + res.status + "\n" + txt)
		console.error("API Error:" + res.status + "\n" + txt)
		let e = new Error(txt)
		throw e
	}

	return txt
}

async function InvokeOk(path: string, method: HTMLMethod, body?: any): Promise<boolean> {

	//@ts-ignore
	let res = await Invoke(path, method, body)

	let txt = await res.text()
	if (res.status >= 400) {
		console.error("API Error:" + res.status + "\n" + txt)
		return false
	}
	return true
}

`, api.BasePath))
	b.WriteString("//#endregion\n\n")
	b.WriteString("//#region Types\n")
	for k, v := range api.Types {
		if v.Desc != "" {
			b.WriteString(fmt.Sprintf("/**\n%s*/\n", v.Desc))
		}

		b.WriteString(fmt.Sprintf("export interface %s {\n", k))
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
				fm, ok := tstypemapper[f.Mapkey]
				if !ok {
					fm = f.Mapkey
				}
				fv, ok := tstypemapper[f.Mapval]
				if !ok {
					fv = f.Mapval
				}
				ftype = "{[s:" + fm + "]:" + fv + "}"
			}

			if f.Desc != "" {
				b.WriteString(fmt.Sprintf("\t/**\n%s*/\n", f.Desc))
			}
			b.WriteString(fmt.Sprintf("\t%s ?: %s\n", strings.ToLower(kf), ftype))
		}

		b.WriteString(fmt.Sprintf("}\n\n"))
	}
	b.WriteString("//#endregion\n\n")
	b.WriteString("//#region Methods\n")
	for k, m := range api.Methods {
		if m.Desc != "" {
			b.WriteString(fmt.Sprintf("/**\n%s*/\n", m.Desc))
		}

		b.WriteString(fmt.Sprintf("export async function %s(req:%s):Promise<%s>{\n", k, _typeName(m.ReqType), _typeName(m.ResType)))
		b.WriteString(fmt.Sprintf("\treturn InvokeJSON(\"%s\",\"%s\",req)\n", m.Path, m.Verb))
		b.WriteString(fmt.Sprintf("}\n\n"))
		//}

	}
	b.WriteString("//#endregion\n")

	err := ioutil.WriteFile(f, b.Bytes(), 0600)
	return err
}
