# go_sdk

## log

### config file

```toml
[log]
Stdout=false # print to terminal
Dir="/tmp/" # log file path
[log.agent]
TaskID="log_test"  # elk id
Addr="127.0.0.1:9977"
proto="udp" # only support udp
```