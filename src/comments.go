package main

import (
	"go/ast"
	"regexp"
	"strings"
)

func manageComments(a *ast.CommentGroup) map[string]string {
	ret := make(map[string]string)

	if a == nil {
		return ret
	}

	var myExp = regexp.MustCompile(`\/?\*?@(?P<k>.*)\s*:\s*(?P<v>.*)\*?\/?`)

	var myExp2 = regexp.MustCompile(`\/?\*?@(?P<k>.*?)\*?\/?$`)

	for _, v := range a.List {
		lines := strings.Split(v.Text, "\n")
		for _, l := range lines {
			m := myExp.FindStringSubmatch(l)
			if m != nil {
				ret[m[1]] = m[2]
			} else {
				m := myExp2.FindStringSubmatch(l)
				if m != nil {
					ret[m[1]] = "OK"
				}
			}

		}
	}
	return ret
}

func manageCommentsGroups(as []*ast.CommentGroup) map[string]string {
	ret := make(map[string]string)

	if as == nil {
		return ret
	}

	for _, a := range as {
		pret := manageComments(a)
		for k, v := range pret {
			ret[k] = v
		}

	}
	return ret
}
