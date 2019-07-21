package tpl

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"os"
	"testing"
	"text/template"
)

func TestGenSvrEntry(t *testing.T) {

	//asset := new(parser.ServerDescriptor)
	//asset.Protocol = "gorpc"

	asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")
	//fmt.Printf("%#v\n", asset)

	tpl_entry, _ := template.New("svr_main.go.tpl").ParseFiles("../asset/go/svr_main.go.tpl")
	tpl_entry.Execute(os.Stdout, *asset)
}

func TestGenSvrInit(t *testing.T) {
	//asset := new(parser.ServerDescriptor)
	//asset.RPC = []parser.ServerRPCDescriptor{}
	//asset.RPC = append(asset.RPC, parser.ServerRPCDescriptor{Name: "BuyApple", Cmd: "BuyApple"})
	//asset.RPC = append(asset.RPC, parser.ServerRPCDescriptor{Name: "SellApple", Cmd: "SellApple"})

	asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")

	tpl_entry, _ := template.New("exec_init.go.tpl").ParseFiles("../asset/go/exec_init.go.tpl")
	tpl_entry.Execute(os.Stdout, asset)
}

func TestGenSvrImpl(t *testing.T) {
	/*
		asset := new(parser.ServerDescriptor)
		asset.PackageName = "test_gorpc"
		asset.Protocol = "gorpc"
		asset.RPC = []parser.ServerRPCDescriptor{}
		asset.RPC = append(asset.RPC, parser.ServerRPCDescriptor{"BuyApple", "BuyApple", "BuyAppleReq", "BuyAppleRsp"})
		asset.RPC = append(asset.RPC, parser.ServerRPCDescriptor{"SellApple", "SellApple", "SellAppleReq", "SellAppleRsp"})
	*/

	asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")
	fmt.Printf("%#v\n", asset)

	tpl_entry, _ := template.New("exec_impl.go.tpl").ParseFiles("../asset/go/exec_impl.go.tpl")
	tpl_entry.Execute(os.Stdout, asset)
}

func TestGenSvrCore(t *testing.T) {

	/*
		asset := new(parser.ServerDescriptor)
		asset.PackageName = "test_gorpc"
		asset.Protocol = "gorpc"
		asset.RPC = []parser.ServerRPCDescriptor{}
		asset.RPC = append(asset.RPC, parser.ServerRPCDescriptor{"BuyApple", "BuyApple", "BuyAppleReq", "BuyAppleRsp"})
		asset.RPC = append(asset.RPC, parser.ServerRPCDescriptor{"SellApple", "SellApple", "SellAppleReq", "SellAppleRsp"})
	*/

	asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")

	tpl_entry, _ := template.New("exec.go.tpl").ParseFiles("../asset/go/exec.go.tpl")
	tpl_entry.Execute(os.Stdout, asset)
}

func TestGenSvrRpc(t *testing.T) {
	//asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")
	//asset, _ := parser.ParseProtoFile("../parser/test.proto", "swan")
	asset, _ := parser.ParseProtoFile("../parser/test.proto", "chick")
	tpl_entry, _ := template.New("rpc.go.tpl").ParseFiles("../asset/go/rpc.go.tpl")
	tpl_entry.Execute(os.Stdout, asset)
}

func TestGenSvrClient(t *testing.T) {
	//asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")
	//asset, _ := parser.ParseProtoFile("../parser/test.proto", "swan")
	asset, _ := parser.ParseProtoFile("../parser/test.proto", "chick")
	tpl_entry, _ := template.New("client.go.tpl").ParseFiles("../asset/go/client.go.tpl")
	tpl_entry.Execute(os.Stdout, asset)
}

func TestGenREADME(t *testing.T) {

	/*
		asset := new(parser.ServerDescriptor)
		asset.ServerName = "hello_svr"
		asset.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	*/

	asset, _ := parser.ParseProtoFile("../parser/test.proto", "gorpc")

	tpl_readme, err := template.New("README.md.tpl").ParseFiles("../asset/go/README.md.tpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	tpl_readme.Execute(os.Stdout, asset)
}
