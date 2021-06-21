# QuickStart

## Installation

### Install 'go'

Firstly, download and install go, see https://golang.org/doc/install.

### Install 'gorpc' tool

'gorpc' is an utility that helps generate server template, RPC Stub, API documentations. It can also help do API testing and pressure testing, etc.

Run following command to install this tool:

```go
go get -u github.com/hitzhangjie/gorpc-cli/gorpc
```

Please make sure you have added `$GOBIN` into your `$PATH` environment variable.

## Write your 1st gorpc server

### define interfaces by protocolbuffers

Taking following protofile as an example:

***file: greeter.proto***

```
syntax = "proto3";
package app;

option go_package = "github.com/examples/helloworld"
    
message Request {};
  
message Response {};
  
service greeter {
    rpc SayHello(Request)  returns(Response);
}
```

Run tool `gorpc` with subcmd `generate` to create template project:

```bash
gorpc generate -protocol=gorpc -protofile=greeter.proto -httpon
```

`gorpc` will create following project template:

```bash
greeter
  |- conf
      |- service.ini
  |- src
      |- main.go
      |- service
          |- service.go
          |- sayHello.go
      |- proto
          |- app
            |- greeter.proto
            |- greeter.pb.go
```

***file: src/service/service.go***

```go
import "proto/app"

type service struct {}

func (s *service) SayHello(ctx context.Context, req app.Request) (rsp *app.Response, err error) {
	
	return
}
```

***file: src/main.go***

```go
import (
	"gorpc/codec/gorpc"
	"gorpc/codec/http"
)


func main() {
    service := gorpc.NewService("gorpc/app.greeter").Version("1.0.0")
    service.Serve(&service.Service{})
}
```

***file: conf/service.ini***

```
[gorpc-service]
tcp.port = 10000
udp.port = 10000

[http-service]
http.port = 8080
http.prefix = /cgi-bin/service/
```

Build and test: `make test`, this will launch a process:
- listen on 10000/tcp and 10001/udp, encoding/decoding via gorpc protocol
- listen on 8080/tcp, encoding/decoding via http protocol
- besides, you can define your business protocol, if needed.

This feature sometimes is very useful when you want to listen both tcp/udp port,
or you want to expose some api interface via another protocol, for example,
to be compatible with old protocol.

Maybe we will don't use it very often, but supporting multiple ports in one process
is still a needed feature in mentioned cases. And this cases are indeed existed. So
I decide to reserve this ability.

## Learn More

Please read the following documentations to learn more.