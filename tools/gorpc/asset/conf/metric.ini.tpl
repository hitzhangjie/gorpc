[{{.ServerName}}]
{{ range .RPC}}
#服务接口-{{.Cmd}}
metric.{{.Cmd}}.timecost10=0                #接口{{.Cmd}}延时10ms
metric.{{.Cmd}}.timecost50=0                #接口{{.Cmd}}延时50ms
metric.{{.Cmd}}.timecost100=0               #接口{{.Cmd}}延时100ms
metric.{{.Cmd}}.timecost300=0               #接口{{.Cmd}}延时300ms
metric.{{.Cmd}}.timecost500=0               #接口{{.Cmd}}延时500ms
metric.{{.Cmd}}.timecost1000=0              #接口{{.Cmd}}延时1000ms
metric.{{.Cmd}}.timecost2000=0              #接口{{.Cmd}}延时2000ms
metric.{{.Cmd}}.timecost3000=0              #接口{{.Cmd}}延时3000ms
metric.{{.Cmd}}.timecostover3000=0          #接口{{.Cmd}}延时>3000ms

{{ end}}