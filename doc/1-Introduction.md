# Introduction

# What's go-rpc ?

`go-rpc` is a rpc framework, it aims to reduce develop and maintenance complexity. 

# Why do we develop go-rpc ?

Compared with standalone application, MicroService Architecture has become more and more popular. Though MicroService Architecture has many advantages over standalone application, it also has many challenges and difficulties to solve, for example:

- naming service, cornerstone to implement remote procedure call, load balancer, scalability, etc.
- load balancer, cornerstone to implement high concurrency, fault torlerance, high availability, etc.
- remote logging, cornerstone to accumulate logging and restore request-process-response occassion.
- monitor system, cornerstone to accumulate, visualize, report system or business events, etc.
- ...

there're many more challenges to solve, it's not that easy as term `micro` means. `go-rpc` is developed to solve this challenges, then reduce develop and maintenance complexity.

`go-rpc` will make you happy with programming again.

## Why not grpc ?

grpc is built upon http/2, we have many services run on tcp, udp and http. Though grpc is good, it doesn't meet our needs. I think many corporations has same occassions as us. So we develop another rpc framework instead of grpc.

go-rpc will support tcp, udp, http/1.x, http/2 and grpc, so you can use go-rpc to handle these different occasions.

## Another Purpose

Also, we want to practise and verify my thinking in software architecture design, programming techniques, etc.

