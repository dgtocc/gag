package main

import (
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"sort"
	"strings"
)

var api API

var knownMethods map[string]bool = make(map[string]bool)
var httpMapper map[string]map[string]string = make(map[string]map[string]string)
var packageName string = "main"

func llog(s string, p ...interface{}) {
	fmt.Printf(s+"\n", p...)
}

func addStruct(a *ast.GenDecl) {
	md := manageComments(a.Doc)
	if md["API"] == "" {
		return
	}
	tp := APIType{
		Name:   "",
		Desc:   "",
		Fields: make(map[string]*APIField),
		Col:    md["COL"],
	}
	tp.Name = a.Specs[0].(*ast.TypeSpec).Name.Name
	llog("Adding type: %s => %#v", tp.Name, md)
	for _, v := range a.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {

		tp.Fields[v.Names[0].Name] = &APIField{}
		tp.Fields[v.Names[0].Name].Tags = make(map[string]APIFieldTag)

		switch x := v.Type.(type) {

		case *ast.Ident:
			tp.Fields[v.Names[0].Name].Type = x.Name
		case *ast.ArrayType:
			switch z := x.Elt.(type) {
			case *ast.Ident:
				tp.Fields[v.Names[0].Name].Type = z.Name
				tp.Fields[v.Names[0].Name].Array = true
			case *ast.InterfaceType:
				tp.Fields[v.Names[0].Name].Type = "interface{}"
				tp.Fields[v.Names[0].Name].Array = true
			case *ast.SelectorExpr:
				api.UsedImportsTypes[z.X.(*ast.Ident).Name] = api.Imports[z.X.(*ast.Ident).Name]
				tp.Fields[v.Names[0].Name].Type = z.X.(*ast.Ident).Name + "." + z.Sel.Name
				tp.Fields[v.Names[0].Name].Array = true
			}

		case *ast.StarExpr:
			switch y := x.X.(type) {
			case *ast.Ident:
				tp.Fields[v.Names[0].Name].Type = y.Name
			case *ast.SelectorExpr:
				switch z := y.X.(type) {
				case *ast.Ident:
					api.UsedImportsTypes[z.Name] = api.Imports[z.Name]
					tp.Fields[v.Names[0].Name].Type = z.Name + "." + y.Sel.Name
				}
			}

		case *ast.InterfaceType:
			tp.Fields[v.Names[0].Name].Type = "interface{}"

		case *ast.SelectorExpr:

			switch z := x.X.(type) {
			case *ast.Ident:
				api.UsedImportsTypes[z.Name] = api.Imports[z.Name]
				tp.Fields[v.Names[0].Name].Type = z.Name + "." + x.Sel.Name
			}

		case *ast.MapType:

			switch z := x.Value.(type) {
			case *ast.Ident:
				tp.Fields[v.Names[0].Name].Type = ""
				tp.Fields[v.Names[0].Name].Mapkey = x.Key.(*ast.Ident).Name
				tp.Fields[v.Names[0].Name].Mapval = z.Name
				tp.Fields[v.Names[0].Name].Map = true
			case *ast.InterfaceType:
				tp.Fields[v.Names[0].Name].Type = "interface{}"
				tp.Fields[v.Names[0].Name].Array = true
			}

		default:
			log.Printf("%#v", x)
		}
		if v.Tag != nil {
			tgstr := strings.ReplaceAll(v.Tag.Value, "`", "")
			tg, err := structtag.Parse(tgstr)
			if err != nil {
				panic(err)
			}
			for _, tgv := range tg.Keys() {
				atg, err := tg.Get(tgv)
				if err != nil {
					panic(err)
				}
				tp.Fields[v.Names[0].Name].Tags[tgv] = APIFieldTag{
					Key:  tgv,
					Name: atg.Name,
					Opts: atg.Options,
				}
			}
			log.Printf("#%v", tg)
		}

	}
	api.Types[tp.Name] = &tp
}

func addFunction(a *ast.FuncDecl) {
	md := manageComments(a.Doc)

	if md["API"] == "" {
		return
	}

	llog("Adding Fuction: %s => %#v", a.Name, md)
	reqType := &APIParamType{}
	resType := &APIParamType{}

	if len(a.Type.Params.List) > 1 {
		switch x := a.Type.Params.List[1].Type.(type) {
		case *ast.StarExpr:
			reqType.Ispointer = true
			switch y := x.X.(type) {
			case *ast.Ident:
				reqType.Typename = y.Name
			case *ast.SelectorExpr:
				api.UsedImportsFunctions[y.X.(*ast.Ident).Name] = api.Imports[y.X.(*ast.Ident).Name]
				reqType.Typename = y.X.(*ast.Ident).Name + "." + y.Sel.Name
			}
		case *ast.ArrayType:
			reqType.IsArray = true
			switch y := x.Elt.(type) {
			case *ast.Ident:
				reqType.Typename = y.Name
			case *ast.SelectorExpr:
				api.UsedImportsFunctions[y.X.(*ast.Ident).Name] = api.Imports[y.X.(*ast.Ident).Name]
				reqType.Typename = y.X.(*ast.Ident).Name + "." + y.Sel.Name
			case *ast.StarExpr:
				reqType.Ispointer = true
				switch z := y.X.(type) {
				case *ast.Ident:
					reqType.Typename = z.Name
				case *ast.SelectorExpr:
					api.UsedImportsFunctions[z.X.(*ast.Ident).Name] = api.Imports[z.X.(*ast.Ident).Name]
					reqType.Typename = z.X.(*ast.Ident).Name + "." + z.Sel.Name
				}
			}
		case *ast.Ident:
			reqType.Typename = x.Name
		case *ast.SelectorExpr:
			api.UsedImportsFunctions[x.X.(*ast.Ident).Name] = api.Imports[x.X.(*ast.Ident).Name]
			reqType.Typename = x.X.(*ast.Ident).Name + "." + x.Sel.Name
		}
	}

	if a.Type.Results != nil && len(a.Type.Results.List) > 0 {

		switch x := a.Type.Results.List[0].Type.(type) {
		case *ast.StarExpr:
			resType.Ispointer = true
			switch y := x.X.(type) {
			case *ast.Ident:
				resType.Typename = y.Name
			case *ast.SelectorExpr:
				api.UsedImportsFunctions[y.X.(*ast.Ident).Name] = api.Imports[y.X.(*ast.Ident).Name]
				resType.Typename = y.X.(*ast.Ident).Name + "." + y.Sel.Name
			}
		case *ast.ArrayType:
			resType.IsArray = true
			switch y := x.Elt.(type) {
			case *ast.Ident:
				resType.Typename = y.Name
			case *ast.SelectorExpr:
				api.UsedImportsFunctions[y.X.(*ast.Ident).Name] = api.Imports[y.X.(*ast.Ident).Name]
				resType.Typename = y.X.(*ast.Ident).Name + "." + y.Sel.Name
			case *ast.StarExpr:
				resType.Ispointer = true
				switch z := y.X.(type) {
				case *ast.Ident:
					resType.Typename = z.Name
				case *ast.SelectorExpr:
					api.UsedImportsFunctions[z.X.(*ast.Ident).Name] = api.Imports[z.X.(*ast.Ident).Name]
					resType.Typename = z.X.(*ast.Ident).Name + "." + z.Sel.Name
				}
			}
		case *ast.Ident:
			resType.Typename = x.Name
		case *ast.SelectorExpr:
			api.UsedImportsFunctions[x.X.(*ast.Ident).Name] = api.Imports[x.X.(*ast.Ident).Name]
			resType.Typename = x.X.(*ast.Ident).Name + "." + x.Sel.Name

		}

		if md["RAW"] == "true" {
			reqType.Typename = md["REQ"]
			resType.Typename = md["RES"]
		}

		verb := md["VERB"]
		if verb == "" {
			verb = http.MethodPost
		}
		fn := APIMethod{
			Name:    a.Name.Name,
			Desc:    a.Name.Name,
			Verb:    verb,
			Path:    md["PATH"],
			Perm:    md["PERM"],
			ReqType: reqType,
			ResType: resType,
		}

		if fn.Path == "" {
			fn.Path = "/" + strings.Replace(strings.ToLower(a.Name.Name), "_", "/", -1)
		}

		api.Methods[a.Name.Name] = &fn
	}

}

func Load(src string) error {

	api.Types = (make(map[string]*APIType))
	api.Methods = (make(map[string]*APIMethod))
	api.Imports = make(map[string]string)
	api.UsedImportsTypes = make(map[string]string)
	api.UsedImportsFunctions = make(map[string]string)
	api.Paths = make(map[string]*APIPath)
	api.SortedPaths = make([]*APIPath, 0)

	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseDir(fset, src, nil, parser.ParseComments)

	if err != nil {
		return err
	}

	for _, v := range f {

		llog("Loading Package: %s", v.Name)
		// Print the AST.
		ast.Inspect(v, func(n ast.Node) bool {

			switch x := n.(type) {

			case *ast.GenDecl:
				if x.Tok == token.TYPE {
					addStruct(x)
				} else {
					return true
				}
			case *ast.ImportSpec:
				var impkey = ""
				var impval = ""
				impval = strings.Replace(x.Path.Value, "\"", "", -1)
				impval = strings.Replace(impval, "'", "", -1)
				if x.Name != nil && x.Name.Name != "" {
					impkey = x.Name.Name
				} else {
					parts := strings.Split(impval, "/")
					impkey = parts[len(parts)-1]
				}
				api.Imports[impkey] = impval
			case *ast.File:
				manageCommentsGroups(x.Comments)
			case *ast.FuncDecl:
				addFunction(x)
				llog("Adding fn: %s", x.Name)
			case *ast.ValueSpec:
				if x.Names[0].Name == "BASEPATH" {
					api.BasePath = strings.Replace(x.Values[0].(*ast.BasicLit).Value, "\"", "", -1)
				}
				if x.Names[0].Name == "NAMESPACE" {
					api.Namespace = strings.Replace(x.Values[0].(*ast.BasicLit).Value, "\"", "", -1)
				}
				log.Printf("%#v", x)
			case *ast.Package:
				packageName = x.Name
			default:
				//log.Printf("%#v", x)
				return true
			}

			return true
		})
	}

	for k, v := range api.Methods {
		path, ok := api.Paths[v.Path]
		if !ok {
			path = &APIPath{
				Path:        v.Path,
				MapVerbs:    make(map[string]*APIVerb),
				SortedVerbs: make([]*APIVerb, 0),
			}
			api.Paths[v.Path] = path
		}
		path.MapVerbs[v.Verb] = &APIVerb{
			Verb:   v.Verb,
			Method: v,
		}
		pathmap, ok := httpMapper[v.Path]
		if !ok {
			httpMapper[v.Path] = make(map[string]string)
			pathmap = httpMapper[v.Path]
		}
		pathmap[v.Verb] = k
	}

	pathNames := make([]string, 0)
	for k, v := range api.Paths {
		verbs := make([]string, 0)
		for k, _ := range v.MapVerbs {
			verbs = append(verbs, k)
		}
		sort.Strings(verbs)

		for _, sv := range verbs {
			v.SortedVerbs = append(v.SortedVerbs, v.MapVerbs[sv])
		}
		pathNames = append(pathNames, k)
	}
	sort.Strings(pathNames)
	for _, p := range pathNames {
		api.SortedPaths = append(api.SortedPaths, api.Paths[p])
	}

	api.Namespace = packageName

	return nil
}
