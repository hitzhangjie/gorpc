{{- $svrName := (index .Services 0).Name -}}
[zipkin]
enabled = false                                             #是否启用zipkin trace
service.name = {{$svrName}}                                 #当前服务名称(span endpoint)
service.addr = *:8000                                       #当前服务地址(span endpoint)
collector.addr = http://a.b.c.d:8080/api/v1/spans           #zipkin collector接口地址
traceId128bits = true                                       #是否启用128bits traceId

[jaeger]
enabled = false                                             #是否启用jaeger trace
