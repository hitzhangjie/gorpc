# Tools

**gorpc**是与gorpc框架配套的一个命令行工具，辅助开发人员开发，主要包含如下功能：

1. 使用 Google Protobuf 作为 IDL (Interface Descriptor Language);
2. gorpc <create>, 指定pb文件，快速生成对应的服务模板、rpc相关client stub、*.pb.go等等;
3. gorpc <update>, 指定pb文件，更新生成的服务模板、rpc相关client stub、*.pb.go等等(开发中)；
4. 其他能力，欢迎提issue，鼓励大家共建。
    - 模板文件，可以在模板目录下任意添加，gorpc会遍历、处理每一个模板文件；
    - 生成服务目录结构与模板中目录结构保持一致；
    - 可以自定义 `gorpc <subcmd>` 实现更多能力， 参考 `gorpc create`, ***file:cmds/create.go***

`gorpc` 实现方式类似于 `go <subcmd>`，subcmd是可插拔的，借此可以来扩展 `gorpc` 的功能。后续也有计划增加其他能力，如快速pprof、console等能力。

## 使用 "Google Protobuf" 作为 IDL

***Google Protobuf*** 是Google开发的具备自描述能力的一种消息格式，与语言无关、平台无关、协议可扩展，应用比较广泛。为了叙述方便，以下简称pb。

pb自身具备的一些特性，使他非常适合用作 IDL (Interface Descriptor Language) 用来指导一些代码生成相关的工作, Google 专门开发了一个针对pb的编译器`protoc`，它能够解析pb文件，并生成与之相关的代码。

两年前，我写过一篇文章详细介绍了 `protoc` 及其插件 (如 `protoc-gen-go`) 之间是如何协作用来生成代码的，如果你对此感兴趣可以读一下这篇文章：[Protoc及其插件工作原理分析](https://hitzhangjie.github.io/2017/05/23/Protoc%E5%8F%8A%E6%8F%92%E4%BB%B6%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86%E5%88%86%E6%9E%90(%E7%B2%BE%E5%8D%8E%E7%89%88).html)。

由于pb语法简单，可读性强，可以较为完整描述一个微服务所需的元信息，`gorpc` 也使用pb文件作为IDL，用来指导服务代码、rpc client stub，甚至是启动配置文件等的生成工作，能一定程度上够简化前期工程搭建的工作。

## 使用 "Go Template" 配置工程模板

`protoc --cpp_out`, `protoc --java_out`, CC++、Java开发中使用过pb的同学，常用上述命令来生成pb文件对应的代码 `*.pb.cc, *.pb.h`, `*.java`。在 pb编译器 `protoc` 中默认内置了某些语言的处理能力，不需要额外的 `protoc plugin` 来支持，但也有些语言的代码生成是没有内置在protoc里面的，如go语言对应的 `protoc-gen-go` 就是单独开发的。此外，如果想自定义代码生成，如支持 `--gorpc_out=`，也需要自行开发 `protoc-gen-gorpc`。

### 概括protoc及其插件工作方式

以 `protoc -go_out=. greeter.proto`为例，介绍下protoc及其插件工作方式。

当protoc执行时，它完成对 `greeter.proto` 文件的解析提取出pb描述信息，并构造一个 `FileDesciptorProto` 对象，该对象包含了greeter.proto文件中的一切必要描述信息。 之后，protoc构造一个代码生成请求 `CodeGenerateRequest`， 该请求中包含了pb文件对应的 `FileDescriptorProto` 对象，然后protoc创建一个子进程启动程序 `protoc-gen-go`，彼此之间通过`pipe`进行通信，protoc将CodeGenerateRequest对象发送给protoc-gen-go
，然后protoc-gen-go开始执行代码生成任务。protoc-gen-go并不直接在本地生成代码，而是将生成的代码内容填充到`CodeGenerateResponse`返回给父进程protoc，由protoc完成最终的代码生成任务。

这就是 `protoc` 及其插件 `protoc-gen-go` 二者的协作方式.

### 选择哪种代码生成方式

本次框架治理，涉及到多语言，包含Go、Java、CC++、NodeJS等，主要有如下考虑：
- 多语言都各自实现一个插件 `protoc-gen-$lang` 涉及到大量重复工作，没有必要，该方案不可取；
- 各语言开发一个共同的子插件 `protoc-gen-gorpc`, 在此基础上扩展子插件(如`plugins=+go`)支持多语言
   代码生成工具往往通过generator g, g.P(...)生成代码，由于要生成的文件、代码数量较多，该中方式调整、维护起来极为不便；
   各语言自定义代码模板，protoc-gen-gorpc内部通过模板引擎处理，将输出内容返回给protoc，这种方式似乎比前一种好一点；
- protoc处理pb文件比大家预想的要复杂一些，尤其是涉及到pb import及指定了其他fileOption（如go_package, java_package, java_outer_classname等）的时候，it's much harder than you think. 如果只是实现protoc插件，那么用户将自己处理这些逻辑比如指定import的pb文件对应的package，`protoc --go_out=Ma/a.proto=aaa`，我相信大部分开发者对protoc掌握的没有这么清楚，暴露这些逻辑只会徒增复杂性；
- 后期业务开发中，可能希望集成mock测试、monitor批量申请、协议管理等能力，如果牵扯到能力类型众多，可能要多个命令行工具；

所以最终选择了这样的实现方式：
- 统一实现一个命令行程序 `gorpc`，其支持自命令`gorpc <subcmd>`，通过子命令来扩展其功能；
- 借助第三方pb解析库，完成pb文件的解析，并将pb描述信息存储到File\Service\Method等层级的Descriptor对象中导出；
- 各语言根据自身需要，自行定制 `go template` 文件，并存放到 `${INSTALL}/asset_${lang}` 目录下；
- gorpc根据命令行参数 `-lang=go` 及配置文件定位到go模板对应的模板目录，并对其下的模板文件逐一处理；

## 如何安装 `gorpc` 命令行工具

### 现阶段安装方式

由于 `gorpc` 除了可执行程序本很，还依赖不同语言的模板文件、gorpc本身配置文件中的信息，`go install` 只安装gorpc命令行工具是无法正常工作的。

现阶段的安装方式：

```bash
git clone https://github.com/hitzhangjie/gorpc-framework
cd gorpc/tools/gorpc
make && make install
```

### 支持 `go install` 安装

后续回考虑通过 `go install` 完成安装，可能采用的方式是：
- 借助go-bindata or go.rice之类的工具将资源文件打包成go文件；
- gorpc链接上述go文件，并在首次运行时将资源文件生成到本地，以支持用户自定义；

由于时间、人力原因，当前还没有支持，如果您有意愿支持，欢迎PR。

## 如何使用 `gorpc` 命令行工具

您可以运行 `gorpc` 或 `gorpc help` 来查看简易帮助信息，会显示各个subcmd对应的使用示例：

```bash
$ gorpc help

how to display help:
        gorpc help

how to create project:
        gorpc create -protodir=. -protofile=*.proto -protocol=whisper -httpon=false
        gorpc create -protofile=*.proto -protocol=gorpc

how to update project:
        gorpc update -protodir=. -protofile=*.proto -protocol=whisper
        gorpc update -protofile=*.proto -protocol=gorpc
```

例如：

```
curl https://git.code.oa.com/gorpc-go/gorpc-go/tree/master/examples/helloworld/helloworld.proto
gorpc create -protofile=helloworld.proto
cd Greeter
go build -v
```

如果遇到问题，可以参考 [HowToBuild](https://git.code.oa.com/gorpc-go/gorpc-go/wikis/Howto/HowToBuild/)

您可以运行 `gorpc help -v` 来查看详细帮助信息，会显示各个subcmd对应的各个选项的信息：

```bash
gorpc help -v

gorpc <cmd> <options>: 

global options:
	-h display this help
	-v display verbose info

gorpc create:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, whisper, or customized ones
	-httpon, enable http mode
	-g, generate code structure conforming to global gopath

gorpc update:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, whisper, or customized ones
```

## 自定义服务模板文件

1. 前面已经提到 `asset_${lang}` 下面的模板文件可以任意添加、删除、修改，gorpc会遍历目录下每个file entry并处理
   - 如果file entry是一个目录, 在输出文件中创建该目录
   - 如果file entry是一个模板文件，执行go模板引擎处理，并在输出文件夹中创建该文件，保留原有的相对路径
   
2. go模板文件中可以使用的一些模板参数信息

   导出给go模板引擎的顶层对象是`FileDescriptor`，结合下面的定义您可以访问pb文件中定义的内容。如可以在模板文件中通过`{{.PackageName}}`来引用`FileDescriptor.PackageName`的值，go template非常简单、灵活，您可以详细阅读相关参考手册，也可以参考已经提供的代码模板`install/asset_go/`来学习如何使用。

    ```go
    // FileDescriptor 文件作用域相关的描述信息
    type FileDescriptor struct {
       PackageName string                 // pb包名称
       Imports     []string               // 跟pb文件中import对应的golang import路径
       FileOptions map[string]interface{} // fileoptions
       Services    []*ServiceDescriptor   // 支持多service
    }
   
    // ServiceDescriptor service作用域相关的描述信息
    type ServiceDescriptor struct {
       Name string           // 服务名称
       RPC  []*RPCDescriptor // rpc接口定义
    }
    
    // RPCDescriptor rpc作用域相关的描述信息
    //
    // RequestType由于涉及到
    type RPCDescriptor struct {
       Name              string // RPC方法名
       Cmd               string // RPC命令字
       FullyQualifiedCmd string // 完整的RPC命令字，用于ServiceDesc、client请求时命令字
       RequestType       string // RPC请求消息类型，包含package，比如package_a.TypeA
       ResponseType      string // RPC响应消息类型，包含package，比如package_b.TypeB
       LeadingComments   string // RPC前置注释信息
       TrailingComments  string // RPC后置注释信息
    }
    ```
   
3. 也提供了为数不多的funcmap函数，供模板中使用
   - title: `{{hello | title}}` ==> `Hello`
   - simplify: `{{simplify helloworld.GreeterServer helloworld}}` ==> `GreeterServer`
   - splitList `{{split "$" "hello$world"}}` ==> `[hello world]`
   - last `{{last (split "/" "git.code.oa.com/abc/def")}}` ==> `def`

## 已知问题

非常希望有意愿的同学一同来改进 `gorpc`工具，提高大家开发效率。

提issue的建议：
- 提新issue之前，请先大致检索是否已经存在类似问题的issue，请不要重复提相同issue；
- 如果您能判断所描述的问题不属于同一个问题、bug、特性，请拆成多个issue进行描述；
- 看到issue是好事，但是我更希望看到issue+PR。

欢迎大家踊跃反馈问题，提PR改进。

