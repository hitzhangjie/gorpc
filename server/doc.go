// Package server provides ability to function as a server.
//
// In gorpc framework, terms `server` and `service` are used:
// - Service, a server instance is actually a server process.
// - Servicce, a service is a provides a kind of ability.
//	 Some relevant APIs are grouped together to form a service.
//
// For example, we have following Google Protobuf file:
//
// file helloworld.proto
// -----------------------------------------------------------
// package helloworld;
//
// message SingleHelloReq {}
// message SingleHelloRsp {}
//
// message MultiHelloReq {}
// message MultiHelloRsp {}
//
// message ByeReq{}
// message ByeRsp{}
//
// service Hello {
//   rpc SayHello(HelloReq) returns(HelloRsp);
//   rpc SayHelloToMulti(MultiHelloReq)  returns(MultiHelloRsp);
// }
//
// service Bye {
//   rpc  SayBye(ByeReq) returns(ByeRsp);
// }
// -----------------------------------------------------------
//
// Well, we can create a server that provides 2 service:
// the HelloService, for SayHello and SayHelloToMulti APIs;
// the ByeService, for SayBye API;
package server
