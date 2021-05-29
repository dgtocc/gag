package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func typeToJson(tn string) interface{} {

	rootmap := make(map[string]interface{})

	td, ok := api.Types[tn]
	if !ok {
		if tn == "string" {
			return "A STRING VALUE"
		}
		if tn == "bool" {
			return true
		}
		if strings.HasPrefix(strings.ToLower(tn), "int") {
			return 123456
		}
		if strings.HasPrefix(strings.ToLower(tn), "float") {
			return 123.456
		}
	}

	for k, v := range td.Fields {
		tg, ok := v.Tags["json"]
		fname := strings.ToLower(k)
		if ok {
			fname = tg.Name
		}
		if v.Map {
			submapval := typeToJson(v.Mapval)
			submap := make(map[string]interface{})
			submap["a"] = submapval
			submap["b"] = submapval
			submap["c"] = submapval
			rootmap[fname] = submap

		} else {
			submapval := typeToJson(v.Type)
			if v.Array {
				submap := make([]interface{}, 0)
				submap = append(submap, submapval)
				submap = append(submap, submapval)
				submap = append(submap, submapval)
				rootmap[fname] = submap
			} else {
				rootmap[fname] = submapval
			}
		}
	}
	return rootmap
}

func typeToJsonStr(tn string) string {
	o := typeToJson(tn)
	bs, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func processHttpCallOut(f string) error {

	b := bytes.Buffer{}

	sortedMethods := make([]string, 0)
	for k, _ := range api.Methods {
		sortedMethods = append(sortedMethods, k)
	}
	sort.Strings(sortedMethods)

	for _, k := range sortedMethods {
		m := api.Methods[k]
		tj := typeToJsonStr(m.ReqType.Typename)
		b.WriteString("###\n")
		if m.Desc != "" {
			b.WriteString(fmt.Sprintf("#%s", strings.Replace(m.Desc, "\n", "\n#", -1)))
		}
		b.WriteString(fmt.Sprintf("\n"))
		b.WriteString(fmt.Sprintf(m.Verb + " https://host/basepath" + m.Path + "\n"))
		b.WriteString("Content-Type: application/json\n")
		b.WriteString("Cookie: dc=<MYCOOKIE>\n\n")
		b.WriteString(tj)
		b.WriteString("\n\n")

	}

	err := ioutil.WriteFile(f, b.Bytes(), 0600)
	return err
}
