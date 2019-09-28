package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// Title uppercase the first character of `s`
func Title(s string) string {
	for idx, c := range s {
		return string(unicode.ToUpper(c)) + s[idx+1:]
	}
	return ""
}

func UnTitle(s string) string {
	for idx, c := range s {
		return string(unicode.ToLower(c)) + s[idx+1:]
	}
	return ""
}

// GoFullyQualifiedType convert $pkg.$type to $realpkg.$type, where $realpkg is calculated
// by `package directive` and `go_package` file option.
func GoFullyQualifiedType(pbFullyQualifiedType string, nfd *FileDescriptor) string {

	rtype := pbFullyQualifiedType

	// 替换RequestType/ResponseType的包名
	idx := strings.LastIndex(pbFullyQualifiedType, ".")
	if idx <= 0 {
		panic(fmt.Errorf("invalid type:%s", pbFullyQualifiedType))
	}

	pkg := pbFullyQualifiedType[0:idx]
	typ := pbFullyQualifiedType[idx+1:]

	if gopkg, ok := nfd.pkgPkgMappings[pkg]; ok && len(gopkg) != 0 {
		rtype = PBGoPackage(gopkg) + "." + typ
	}
	return rtype
}

// PBSimplifyGoType determine whether to use fullyQualifiedPackageName or not,
// if the `fullTypeName` occur in code of `package goPackageName`, `package` part
// should be removed.
func PBSimplifyGoType(fullTypeName string, goPackageName string) string {
	//fmt.Println("fullyQualified:", fullTypeName, "goPackage:", goPackageName)

	idx := strings.LastIndex(fullTypeName, ".")
	if idx <= 0 {
		panic(fmt.Sprintf("invalid fullyQualifiedType:%s", fullTypeName))
	}

	pkg := fullTypeName[0:idx]
	typ := fullTypeName[idx+1:]

	if pkg == goPackageName {
		//fmt.Println("pkg:", pkg, "=", "gopkg:", goPackageName)
		return typ
	}

	//fmt.Println("pkg:", pkg, "!=", "gopkg:", goPackageName)
	return fullTypeName
}

// PBGoPackage convert a.b.c to a_b_c
func PBGoPackage(pkgName string) string {
	var (
		prefix string
		pkg    string
	)
	idx := strings.LastIndex(pkgName, "/")
	if idx < 0 {
		pkg = pkgName
	} else {
		prefix = pkgName[0:idx]
		pkg = pkgName[idx+1:]
	}

	gopkg := strings.Replace(pkg, ".", "_", -1)

	if len(prefix) == 0 {
		return gopkg
	}
	return prefix + "/" + gopkg
}

// PBGoType convert `t` to go style (like a.b.c.hello, it'll be changed to a_b_c.Hello)
func PBGoType(t string) string {

	idx := strings.LastIndex(t, ".")
	if idx <= 0 {
		panic(fmt.Sprintf("fatal error: invalid type:%s", t))
	}

	gopkg := PBGoPackage(t[0:idx])
	msg := t[idx+1:]

	return GoExport(gopkg + "." + msg)
}

// GoExport export go type
func GoExport(typ string) string {
	idx := strings.LastIndex(typ, ".")
	if idx < 0 {
		return strings.Title(typ)
	}
	return typ[0:idx] + "." + strings.Title(typ[idx+1:])
}

// SplitList split string `str` via delimiter `sep` into a list of string
func SplitList(sep, str string) []string {
	return strings.Split(str, sep)
}

// Last returns the last element in `list`
func Last(list []string) string {
	idx := len(list) - 1
	return list[idx]
}

// TrimRight trim right substr starting at `sep`
func TrimRight(sep, str string) string {
	idx := strings.LastIndex(str, sep)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

func HasPrefix(prefix, str string) bool {
	return strings.HasPrefix(str, prefix)
}
