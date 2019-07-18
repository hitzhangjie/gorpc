package parser

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var verbose = flag.Bool("v", true, "verbose")

func TestMain(m *testing.M) {
	flag.Parse()
	ret := m.Run()
	os.Exit(ret)
}

func TestParseProtoFile(t *testing.T) {
	fmt.Println("nrpc test:")
	ParseProtoFile("../test/nrpc/test_nrpc.proto", "nrpc")
	fmt.Println()

	fmt.Println("simplesso test:")
	ParseProtoFile("../test/simplesso/test_simplesso.proto", "simplesso")
	fmt.Println()
}
