#!/bin/sh

# command to learn all server method

## Below command describe all grpc methods
grpcurl -protoset ./script.protoset describe ScriptService


echo "command to create script on server content of file : echo "\\"Hello World"\"
grpcurl -plaintext -protoset ./script.protoset -d '{"script_name":"test", "Content": "IyEvYmluL3NoCgplY2hvICJoZWxsbyB3b3JsZCI="}' localhost:9000 ScriptService/CreateTask
#**************************************************
#output expected
#{
#  "scriptName": "test",
#  "scriptStatus": "Created"
#}
#**************************************************
echo "command to execute task "
grpcurl -plaintext -protoset ./script.protoset -d '{"script_name":"test"}' localhost:9000 ScriptService/ExecuteTask
#**************************************************
#output expected
#{
 # "scriptName": "test",
#  "scriptStatus": "CREATED",
#  "lastRunStatus": "Executed",
#  "output": "hello world\n"
#}
#**************************************************
echo "command to get task  status"
grpcurl -plaintext -protoset ./script.protoset -d '{"script_name":"test"}' localhost:9000 ScriptService/GetTaskStatus
#**************************************************
#output expected
#{
#  "scriptName": "test",
#  "scriptStatus": "CREATED",
#  "lastRunStatus": "Executed"
#}
#**************************************************
echo "command to get task  source"
grpcurl -plaintext -protoset ./script.protoset -d '{"script_name":"test"}' localhost:9000 ScriptService/GetTaskSource
#**************************************************
#output expected
#{
#  "scriptName": "test",
#  "scriptStatus": "CREATED",
#  "content": "#!/bin/sh\n\necho \"hello world\""
#}
#**************************************************