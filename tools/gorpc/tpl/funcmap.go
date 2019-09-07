package tpl

import (
	"git.code.oa.com/trpc-go/trpc-go-cmdline/parser"
	"text/template"
)

var funcMap = template.FuncMap{
	"simplify":   parser.PBSimplifyGoType,
	"gopkg":      parser.PBGoPackage,
	"gotype":     parser.PBGoType,
	"export":     parser.GoExport,
	"gofulltype": parser.GoFullyQualifiedType,
	"title":      parser.Title,
	"trimright":  parser.TrimRight,
	"splitList":  parser.SplitList,
	"last":       parser.Last,
	"hasprefix":  parser.HasPrefix,
}
