package params

import (
	"flag"
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
)

// RepeatedOption support repeated parameters cmdline option
type RepeatedOption []string

func (l *RepeatedOption) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *RepeatedOption) Get() interface{} {
	return *l
}

func (l *RepeatedOption) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func (l *RepeatedOption) Replace(arr *[]string) {
	*l = *arr
}

// LookupFlag lookup a flag value
func LookupFlag(flagSet *flag.FlagSet, pname string, pval interface{}) {

	f := flagSet.Lookup(pname)
	if f == nil {
		return
	}

	v := f.Value.(flag.Getter).Get()
	switch pval.(type) {
	case *string:
		*(pval.(*string)) = v.(string)
	case *bool:
		*(pval.(*bool)) = v.(bool)
	case *int:
		*(pval.(*int)) = v.(int)
	case *RepeatedOption:
		*(pval.(*RepeatedOption)) = v.(RepeatedOption)
	default:
		log.Error("Unknown flag type:[%T]", pval)
	}
}
