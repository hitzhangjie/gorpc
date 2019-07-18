package params

import "strings"

type StringArray []string

func (strArr *StringArray) String() string {
	return strings.Join(*strArr, " ")
}

func (strArr *StringArray) Get() interface{} {
	return *strArr
}

func (strArr *StringArray) Set(value string) error {
	*strArr = append(*strArr, value)
	return nil
}

func (strArr *StringArray) Replace(arr *[]string) {
	*strArr = *arr
}
