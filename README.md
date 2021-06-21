# GoRPC Framework

**GoRPC** is a lightweight RPC framework, **"Simple & Powerful"**. 

GoRPC uses **Google Protobuf** as **IDL (Interface Descriptor Language)**, it supports:

- quickly generate server template, rpc server stub and rpc client stub
- support multiple port and multiple protocols in the same process
- support gorpc, http, and user customized business protocol
- support tcp, udp, unix server and client
- support naming service, default is etcd
- support metrics, default is ouput to local stat file
- support tracing, default is zipkin or jaeger
- support logging, default is local logging

# QuickStart

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

# Efficient Developer Tools
| tools | function | remark |
|:-----:|:--------:|:------:|
|gorpc|quickly generate code ||
|bot|bot of qq, wework, slack||
|log|logging analysis tool||

Actually, gorpc is pluggable like `go <subcmd>` mechanism.

# Go-RPC Arch Design

![arch](https://github.com/hitzhangjie/gorpc/blob/master/docs/arch.png)


# Thanks

![jetbrains.svg](https://github.com/hitzhangjie/gorpc/blob/master/docs/jetbrains.svg)

Special thanks to the [Jetbrains Team](https://www.jetbrains.com/?from=gorpc) for their support.

