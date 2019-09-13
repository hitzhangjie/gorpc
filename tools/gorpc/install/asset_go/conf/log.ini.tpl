#框架内部日志
[log-gorpc_frame]
level = 1                       #日志级别,0:DEBUG,1:INFO,2:WARN,3:ERROR
logwrite = rolling
logFileAndLine = 1
rolling.filename = gorpc_frame.log
rolling.type = size
rolling.filesize = 100m
rolling.lognum = 5

#框架流水日志
[log-gorpc_access]
level = 1                      #日志级别,0:DEBUG,1:INFO,2:WARN,3:ERROR)
logwrite = rolling
logFileAndLine = 0
rolling.filename = gorpc_access.log
rolling.type = daily
rolling.lognum = 5

#服务默认日志
[log-default]
level = 1                     #日志级别,0:DEBUG,1:INFO,2:WARN,3:ERROR)
logwrite = rolling
logFileAndLine = 0
rolling.filename = default.log
rolling.type = size
rolling.filesize = 100m
rolling.lognum = 5