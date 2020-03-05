package gorpc

import "fmt"

// Version gorpc version
type Version struct {
	Major    string
	Minor    string
	Patch    string
	Metadata string
}

// GoRPCVersion is the current version of Delve.
var GoRPCVersion = Version{
	Major:    "0",
	Minor:    "1",
	Patch:    "0",
	Metadata: "dev",
}

func (v Version) String() string {
	ver := fmt.Sprintf("Version: %s.%s.%s", v.Major, v.Minor, v.Patch)
	if v.Metadata != "" {
		ver += "-" + v.Metadata
	}
	return ver
}
