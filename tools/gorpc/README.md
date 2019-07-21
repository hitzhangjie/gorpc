# gorpc

***gorpc*** is an efficient tool to help developers :
- quickly generate or update service template 
- quickly generate service client stub
- manage your protocol buffers, like pull, push

gorpc works like `go <subcmd>`, it's pluggable, which is good to extend its abilities.

## Using "Google Protobuf" as IDL

***Google Protobuf*** is developed by Google, it's a self-descriptive message format.

Using protobuf as IDL (Interface Descriptor Language) is very common, Google also 
provides a protobuf compiler called `protoc`.

Before, I wrote an article which describes the `protoc`, protoc plugins like `proto-gen-go`, 
and `protobuf` internals. If you're interested, please read my article:
[Protoc及其插件工作原理分析(精华版)](https://hitzhangjie.github.io/2017/05/23/Protoc%E5%8F%8A%E6%8F%92%E4%BB%B6%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86%E5%88%86%E6%9E%90(%E7%B2%BE%E5%8D%8E%E7%89%88).html).

## Using "Go Template" as weapon

`protoc --cpp_out`, `protoc --java_out`, we usually use these commands to generate cpp 
header/source files or java files, this works for `cpp`, `java`.

While for many other languages, `protoc` doesn't implement it, such as, go programming 
language. If we want to generate `*.pb.go` like `*.pb.h/*.pb.cc` or `*.pb.java`, we 
should implement a plugin for `go`.

But how ?

Just now, we know protobuf is a self-descriptive message format, when file `*.proto` 
parsed by `protoc`, a `FileDesciptorProto` object will be built, it contains nearly
everything about the `*.proto` we written. If you know little about internals of `protoc`
or `protobuf` itself, please refer to my article metioned above.

when run command `protoc --go_out *.proto`, protoc will read your protofile and parse it,
after that, it build a FileDescriptorProto object, then it will serialize it and search
executable named `protoc-gen-go` in your `PATH` shell env variable. If found, it will
fork a childprocess to run `protoc-gen-go`, and parentprocess `protoc` will create a 
pipe btw itself and childprocess to communicate. `protoc` will send a `CodeGenerateRequest`
to the childprocess `protoc-gen-go` via pipe. This `CodeGenerateRequest` contains 
serialized `FileDescriptorProto`, then `protoc-gen-go` read from pipe and extract it.
`protoc-gen-go` will be responsible for generate source code by `g.P("..")`. This generated
source code info will be responded to `protoc`, `protoc` process will create file and 
write file content (source code).

This is the way `protoc` and `protoc-gen-go` works.

Writing a plugin `protoc-gen-go` or `protoc-gen-gorpc` is really not a good idea for 
generating source code, because it increases the difficulty in maintenance and 
extensibility (for example, generate Java/Python source code). We should parse the
*.proto file once, then using template technology to generate files. If so, all 
programming languages can provide a template, we can quickly generate the service
template, needless to modify the code of `protoc-gen-gorpc`.

So, we'll use some thirdparty protoparsing library to parse *.proto, then use go template
as our weapons to generate service template, client stub, even service configurations.

## How to use gorpc ?

It's user friendly, all subcmds and its options are described in detail.

```bash
$ gorpc help

how to display help:
        gorpc help

how to create project:
        gorpc create -protodir=. -protofile=*.proto -protocol=gorpc -httpon=false
        gorpc create -protofile=*.proto -protocol=gorpc

how to update project:
        gorpc update -protodir=. -protofile=*.proto -protocol=gorpc
        gorpc update -protofile=*.proto -protocol=gorpc
```

or 

```bash
gorpc help -v


gorpc <cmd> <options>: 

global options:
	-h display this help
	-v display verbose info

gorpc create:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, gorpc, chick or swan
	-httpon, enable http mode
	-g, generate code structure conforming to global gopath

gorpc update:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, gorpc, chick or swan
```

