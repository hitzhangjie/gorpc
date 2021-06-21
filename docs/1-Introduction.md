# Introduction

# What's gorpc ?

`gorpc` is an RPC framework, it aims to reduce the complexity of development and maintenance. 

# Why do we develop gorpc ?

Compared with standalone application, MicroService Architecture has become more and more popular. Though MicroService Architecture has many advantages over standalone application, it also has many challenges and difficulties to solve, for example:

- naming service, cornerstone to implement remote procedure call, load balancer, scalability, etc.
- load balancer, cornerstone to implement high concurrency, fault torlerance, high availability, etc.
- remote logging, cornerstone to accumulate logging and restore request-process-response occassion.
- monitor system, cornerstone to accumulate, visualize, report system or business events, etc.
- networking communication, cornerstone to build high-performance backend servers.
- ...

there're many more challenges to solve, it's not that easy as term `micro` means. `gorpc` is developed to solve this challenges, then reduce develop and maintenance complexity.

`gorpc` will make you happy with programming again.

## Why not using grpc, instead ?

grpc is built upon http/2, we have many services run on tcp, udp and http. Though grpc is good, it doesn't meet our needs. I think many corporations has same occassions as us. So we develop another RPC framework instead of grpc.

gorpc will support tcp, udp, http/1.x, http/2, http/3 and grpc, so you can use gorpc to handle these different occasions.

Besides, we want to practise and verify my thinking in software architecture design, programming techniques, etc.

After several years' working for building high-performance and extensible RPC framework, we have learned some experience and skills. Now we want to share this knowledge to others.

That's the background we built gorpc!

## Summary

It's hard to say which RPC framework is better. There're many RPC frameworks, any of which has its own creativity. We choose a proper framework according to many forces, like programming languages, stability, performance, stars, contributors, ecosystem, etc.

We hope `gorpc` could help others understand how to develop an better RPC framework.

Use at your own risk in your productive environment.