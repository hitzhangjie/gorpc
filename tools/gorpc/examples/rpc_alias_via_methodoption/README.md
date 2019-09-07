# generate *.pb.go

生成 ****.pb.go*** 的做法，通过大家习惯的方式，直接借助 ***protoc*** 命令来生成：
```bash
protoc --go_out=. helloworld.proto
```

# generate server stub

通过 [jhump/protoreflect]() 来解析 *.proto 文件，解析后得到 ***FileDescriptorProto*** 对象，该对象中记录了我们关心的 *.proto 描述信息，然后将这些信息填充到go、java、c++等提供的模板中，完成server、client stub的代码生成。

对于server stub，这里需要注意与框架代码的整合：
- 一方面要符合业务开发人员的习惯，使业务开发工作尽可能简单，如只填充业务函数即可；
- 一方面要将业务代码和框架代码进行粘合，组合框架能力、默认插件实现、整合wrapper方法等实现开箱即用；

server要处理的请求命令字或者rpc接口，都已经在*.proto中定义：
- `gorpc` 根据service中定义的rpc，建立映射关系rpcName->Handler；
- `gorpc` 提供方法，使得server能够注册上述映射关系，并在router中根据req找到Handler；

可以参考 `protoc --go_out=plugins=grpc:.` 与 `protoc --go_out=.` 两种方式生成的代码的区别，来参考下grpc是如何建立映射关系的。

## demo

```protobuf
syntax = "proto3";
package helloworld;

// Hello
message HelloRequest {
    string from     = 1;    // say hello from
    string to       = 2;    // say hello to
    string words    = 3;    // hello words
}

message HelloResponse {
    uint32 errcode  = 1;    // error code
    string errmsg   = 2;    // error msg
}

// Bye
message ByeRequest {
    string from     = 1;    // say bye from
    string to       = 2;    // say bye to
    string words    = 3;    // bye words
}

message ByeResponse {
    uint32 errcode  = 1;    // error code
    string errmsg   = 2;    // error msg
}

// service: greeter
service greeter {
    rpc SayHello ( HelloRequest )    returns ( HelloResponse );
    rpc SayBye   ( ByeRequest   )    returns ( ByeResponse   );
}
```

运行gorpc来生成server工程：`gorpc create -protofile=helloworld.proto -protocol=gorpc`，需要生成的内容应该包括：

- helloworld.pb.go
- service interface

```go
import pb "helloworld.proto"

type GreeterServer interface {
    SayHello(ctx context.Context, req *pb.HelloRequest) (rsp *pb.HelloResponse, error)
    SayBye(ctx context.Context, req *pb.ByeRequest) (rsp *pb.ByeResponse, error)
}
```

- service rpcName->rpcMethod register 

```go
func RegisterGreeterServer(s *server.Server, svr GreeterServer) {
    s.RegisterService(&_Greeter_serviceDesc, svr)
}

func _Greeter_SayHello_Handler(svr interface{}, ctx context.Context) error {
    reqCtx := ctx.Value(ctxkey)

    req := new(HelloRequest)
    if err := reqCtx.Decode(req); err != nil {
        return err
    }

    if rsp, err := svr.SayHello(ctx, req); err != nil {
        return err
    }

    reqCtx.rspChan <- rsp
    return nil
}

var _Greeter_serviceDesc = gorpc.ServiceDesc{
    ServiceName: "helloworld.greeter",
    HandlerType: (*GreeterServer)(nil),
    Methods: []gorpc.MethodDesc{
        {
            MethodName: "SayHello",
            Handler: _Greeter_SayHello_Handler,
        },
        {
            MethodName: "SayBye",
            Handler: _Greeter_SayBye_Handler,
        }
    },
    Streams: []gorpc.StreamDesc{},
    MetaData: "helloworld.proto",
}

```

当前server端提供了一个server.WithHandler(....)来封装所有的请求处理、tracing、拦截器、监控、logging等逻辑，这个没问题，Handler内部会请求Dispatcher来完成rpc请求到rpc处理函数的分发。所以这里还需要提供一个Dispatcher？

- 方法1：gorpc可以显示提供一个dispatcher出来，server端使用的时候WithDispatcher(GreeterServer.Dispatcher)就可以；
- 方法2：新增一个方法server.RegisterService(gorpc.ServiceDesc)，server自己注册；
- 方法3：Dispatcher接口提供Add等方法，支持直接注册到dispatcher；
- 方法4：server提供Dispatch(rpcName, rpcMethod)直接进行注册，内部注册到server.Opts.Dispatcher上；

有多种方式，需要综合考虑下，对于后面支持进程多server实例有用处！

# generate client stub

# generate config 

# generate others
