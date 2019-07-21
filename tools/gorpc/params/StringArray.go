package params

import (
	"fmt"
)

type List []string

func (l *List) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *List) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func (l *List) Replace(arr *[]string) {
	*l = *arr
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
