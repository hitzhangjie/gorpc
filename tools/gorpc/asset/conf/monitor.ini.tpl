[{{.ServerName}}]
{{ range .RPC}}
#服务接口-{{.Cmd}}
monitor.{{.Cmd}}.timecost10=0                #接口{{.Cmd}}延时10ms
monitor.{{.Cmd}}.timecost20=0                #接口{{.Cmd}}延时20ms
monitor.{{.Cmd}}.timecost50=0                #接口{{.Cmd}}延时50ms
monitor.{{.Cmd}}.timecost100=0               #接口{{.Cmd}}延时100ms
monitor.{{.Cmd}}.timecost200=0               #接口{{.Cmd}}延时200ms
monitor.{{.Cmd}}.timecost300=0               #接口{{.Cmd}}延时300ms
monitor.{{.Cmd}}.timecost400=0               #接口{{.Cmd}}延时400ms
monitor.{{.Cmd}}.timecost500=0               #接口{{.Cmd}}延时500ms
monitor.{{.Cmd}}.timecost600=0               #接口{{.Cmd}}延时600ms
monitor.{{.Cmd}}.timecost700=0               #接口{{.Cmd}}延时700ms
monitor.{{.Cmd}}.timecost800=0               #接口{{.Cmd}}延时800ms
monitor.{{.Cmd}}.timecost900=0               #接口{{.Cmd}}延时900ms
monitor.{{.Cmd}}.timecost1000=0              #接口{{.Cmd}}延时1000ms
monitor.{{.Cmd}}.timecost2000=0              #接口{{.Cmd}}延时2000ms
monitor.{{.Cmd}}.timecost3000=0              #接口{{.Cmd}}延时3000ms
monitor.{{.Cmd}}.timecostover3000=0          #接口{{.Cmd}}延时>3000ms

{{ end}}