{{- $svrName := (index .Services 0).Name -}}
[service]
name = {{$svrName}}                  #服务名称
limit.max.reqs = 100000              #服务允许最大qps
limit.max.conns = 100000             #允许最大入连接数
tcp.conn.idletime = 300000           #tcp连接空闲关闭时间,5min
workerpool.size = 20000              #worker数量
udp.buffer.size = 4096               #udp接收缓冲大小(B),默认1KB,请注意收发包尺寸
#命令字超时时间设置

[{{.Protocol}}-service]
tcp.port = 8000                      #tcp监听端口
udp.port = 8000                      #udp监听端口

{{- if .HttpOn}}

[http-service]
http.port = 8080                     #监听http端口
http.prefix = /cgi-bin/web           #httpUrl前缀
{{- end}}

[rpc-{{$svrName}}]
addr = ip://127.0.0.1:8000           #rpc调用地址
trans = 3                            #网络传输模式,UDP,TCP_SHORT,TCP_KEEPALIVE,...
timeout = 1000                       #rpc全局默认timeout
#具体到接口的超时时间设置
