[service]
name = {{.ServerName}}               #服务名称
limit.reqs = 100000                  #服务允许最大qps
limit.conns = 100000                 #允许最大入连接数
tcp.conn.idletime = 300000           #tcp连接空闲关闭时间,5min
workerpool.size = 20000              #worker数量
udp.buffer.size = 4096               #udp接收缓冲大小(B),默认1KB,请注意收发包尺寸
{{- range .RPC}}
{{.Cmd}}.cmd.timeout = 5000          #服务接口{{.Cmd}}超时时间(ms)
{{- end}}
env = test

[{{.Protocol}}-service]
tcp.port = 8000                      #tcp监听端口
udp.port = 8000                      #udp监听端口

{{- if .HttpOn}}

[http-service]
http.port = 8080                     #监听http端口
http.prefix = /cgi-bin/web           #httpUrl前缀
{{- end}}

[rpc-{{.ServerName}}]
addr = ip://127.0.0.1:8000           #rpc调用地址
proto = 3                            #网络传输模式,1:UDP,2:TCP_SHORT,3:TCP_KEEPALIVE,4:TCP_FULL_DUPLEX,5:UDP_FULL_DUPLEX,6:UDP_WITHOUT_RECV
timeout = 1000                       #rpc全局默认timeout
{{- range .RPC}}
{{.Cmd}}.timeout = 1000              #rpc-{{.Cmd}}超时时间(ms)
{{- end}}
