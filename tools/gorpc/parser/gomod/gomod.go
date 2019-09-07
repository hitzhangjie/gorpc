package gomod

import (
	"bufio"
	"os"
	"path"
	"strings"
)

func LoadGoMod() (mod string, err error) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := path.Join(d, "go.mod")
	_, err = os.Lstat(p)
	if err != nil {
		return
	}
	fin, err := os.Open(p)
	if err != nil {
		return
	}
	sc := bufio.NewScanner(fin)
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "module ") {
			return strings.Split(l, " ")[1], nil
		}
	}
	return
}
