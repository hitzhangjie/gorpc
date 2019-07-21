package params

import (
	"flag"
	"fmt"
)

func Lookup(pname string, pval interface{}) {

	if f := flag.Lookup(pname); f == nil {

		return

	} else {

		v := f.Value.(flag.Getter).Get()

		switch pval.(type) {
		case *string:
			*pval.(*string) = v.(string)
		case *bool:
			*pval.(*bool) = v.(bool)
		case *int:
			*pval.(*int) = v.(int)
		case *List:
			*pval.(*List) = v.(List)
		default:
			fmt.Println("Unknown flag type:[%T]", pval)
		}
	}
}
