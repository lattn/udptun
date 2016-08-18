# udptun

`udptun`可以监听主机多个端口上的 udp 报文数据,并进行转发。

# 使用方法

在当前工作目录下的`*.json`文件都会被作为一个转发服务的单独配置。

执行命令:
 
 `go run path/to/udptun/cmd/main.go`
 
 # 配置
 
```json
{
  "local_addr": ":25501",
  "target_addr": "192.168.200.1:8877",
  "allowed_ips": [
    "111.203.228.26"
  ],
  "error_msg": "error",
  "buffer_size": 65535
}
```