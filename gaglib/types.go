package gaglib

type API struct {
	BasePath             string                `yaml:"basepath,omitempty" xml:"base_path" json:"base_path,omitempty"`
	Host                 string                `yaml:"host,omitempty" json:"host,omitempty" xml:"host"`
	Types                map[string]*APIType   `yaml:"types,omitempty" json:"types,omitempty" xml:"types"`
	Methods              map[string]*APIMethod `yaml:"methods,omitempty" json:"methods,omitempty" xml:"methods"`
	Namespace            string                `yaml:"namespace" json:"namespace,omitempty" xml:"namespace"`
	Imports              map[string]string     `yaml:"imports" json:"imports,omitempty" xml:"imports"`
	UsedImportsTypes     map[string]string     `yaml:"used_imports_types" json:"used_imports_types,omitempty" xml:"used_imports_types"`
	UsedImportsFunctions map[string]string     `yaml:"used_imports_functions" json:"used_imports_functions,omitempty" xml:"used_imports_functions"`
	SortedPaths          []*APIPath            `yaml:"sorted_paths" json:"sorted_paths,omitempty" xml:"sorted_paths"`
	Paths                map[string]*APIPath   `yaml:"paths" json:"paths,omitempty" xml:"paths"`
}

type APIPath struct {
	Path        string              `yaml:"path" xml:"path"`
	MapVerbs    map[string]*APIVerb `yaml:"map_verbs" xml:"map_verbs"`
	SortedVerbs []*APIVerb          `yaml:"sorted_verbs" xml:"sorted_verbs"`
}

type APIVerb struct {
	Verb   string     `yaml:"verb" xml:"verb"`
	Method *APIMethod `yaml:"method" xml:"method"`
}
type APIFieldTag struct {
	Key  string   `yaml:"key" xml:"key"`
	Name string   `yaml:"name" xml:"name"`
	Opts []string `yaml:"opts" xml:"opts"`
}
type APIField struct {
	Type   string                 `yaml:"type,omitempty" xml:"type"`
	Array  bool                   `yaml:"array,omitempty" xml:"array"`
	Desc   string                 `yaml:"desc,omitempty" xml:"desc"`
	Map    bool                   `yaml:"map,omitempty" xml:"map"`
	Mapkey string                 `yaml:"mapkey,omitempty" xml:"mapkey"`
	Mapval string                 `yaml:"mapval,omitempty" xml:"mapval"`
	Tags   map[string]APIFieldTag `yaml:"tags,omitempty" xml:"tags"`
}

func (a *APIField) String() string {
	if a.Array {
		return "[]" + a.Type
	} else {
		return a.Type
	}
}

type APIType struct {
	Name    string               `yaml:"name,omitempty" xml:"name"`
	Desc    string               `yaml:"desc,omitempty" xml:"desc"`
	Fields  map[string]*APIField `yaml:"fields,omitempty" xml:"fields"`
	Col     string               `yaml:"col,omitempty" xml:"col"`
	TypeDef string               `yaml:"-"`
}

type APIParamType struct {
	Typename  string `xml:"typename"`
	Ispointer bool   `xml:"ispointer"`
	IsArray   bool   `xml:"is_array"`
}

type APIMethod struct {
	Name    string `yaml:"name" xml:"name"`
	Desc    string `yaml:"desc" xml:"desc"`
	Verb    string `yaml:"verb" xml:"verb"`
	Path    string `yaml:"path" xml:"path"`
	Perm    string `yaml:"perm" xml:"perm"`
	ReqType *APIParamType
	ResType *APIParamType
}

func APIParamTypeToString(t *APIParamType) string {
	ret := ""
	if t.IsArray {
		ret = "[]"
		if t.Ispointer {
			ret = ret + "*"
		}
		ret = ret + t.Typename
		return ret
	}
	if t.Ispointer {
		ret = ret + "*"
	}
	ret = ret + t.Typename
	return ret
}

func APIParamTypeDecToString(t *APIParamType) string {
	ret := ""
	if t.IsArray {
		ret = "[]"
		if t.Ispointer {
			ret = ret + "*"
		}
		ret = ret + t.Typename
		return ret
	}
	if t.Ispointer {
		ret = ret + "*"
	}
	ret = ret + t.Typename
	return ret
}

func APIParamTypeUseRef(t *APIParamType) string {
	if t.IsArray || !t.Ispointer {
		return "&"
	}
	return ""
}
