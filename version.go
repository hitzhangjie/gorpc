package gorpc

import "fmt"

// Version gorpc version
type Version struct {
	Major    string
	Minor    string
	Patch    string
	Metadata string
	Build    string
}


var (
	// GoRPCVersion is the current version of Delve.
	GoRPCVersion = Version{
		Major: "0", Minor: "1", Patch: "0", Metadata: "",
		Build: "$Id: ddcc99e754f9c96328ac82c53a4af978f71a7df3 $",
	}
)

func (v Version) String() string {
	ver := fmt.Sprintf("Version: %s.%s.%s", v.Major, v.Minor, v.Patch)
	if v.Metadata != "" {
		ver += "-" + v.Metadata
	}
	return fmt.Sprintf("%s\nBuild: %s", ver, v.Build)
}