package params

import (
	"fmt"
)

type List []string

func (l *List) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *List) Get() interface{} {
	return *l
}

func (l *List) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func (l *List) Replace(arr *[]string) {
	*l = *arr
}
