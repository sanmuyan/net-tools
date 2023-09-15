# 网络工具箱

## 端口扫描

```shell
portscan 192.168.1.0/24 22
Opened: 192.168.1.1:22
Opened: 192.168.1.3:22
Opened: 192.168.1.4:22
Opened: 192.168.1.11:22
Opened: 192.168.1.15:22
```

- `-t` 扫描最大线程，默认值`1`
- `-T` 扫描超时时间（毫秒），默认值`200`

## TCP 测试

```shell
tcpping.exe www.baidu.com 443
Reply from www.baidu.com:443 time=13ms
Reply from www.baidu.com:443 time=9ms
Reply from www.baidu.com:443 time=20ms
Reply from www.baidu.com:443 time=12ms
Success=4, Error=0, Max=20ms, Min=9ms, Avg=13ms
```

- `-c` PING 的次数，默认值`4`
- `-i` 每次PING 的间隔（毫秒），默认值`1`
- `-P` 指定协议，默认值`tcp`，可选项 `tcp|http|read`
- `-T` PING 超时时间

### 协议说明

- `tcp` 建立TCP 连接后即认为成功
- `http` 主动发送一个HTTP `HEAD`请求，并读取HTTP 服务器响应
- `read` 建立连接后，读取一次服务端发送的`hello`包

## 带宽测试

### 服务端

```shell
speedtests
2023/08/01 18:02:25 tcp server runing 0.0.0.0:8080
2023/08/01 18:03:02 tcp download from 192.168.1.2:50916
2023/08/01 18:03:12 tcp download finished in 192.168.1.2:50916
```

- `-p` 服务器监听端口，默认值`8080`
- `-s` 服务器监控地址，默认值`0.0.0.0`
- `-P` 服务器监听协议，默认值`tcp`，可选项`tcp|udp`

### 服务端

```shell
speedtestc 192.168.1.1:8080
2023/08/01 18:03:03 tcp download testing to 192.168.1.1:8080
2023/08/01 18:03:13 finished speed: 897Mbps/s
```

- `-m` 测试模式，默认值`download`，可选项`download|upload`
- `-t` 测试时间（秒），默认值`10`
- `-P` 测试的协议，默认值`tcp`，可选项`tcp|udp`
- `-T` 测试并发连接，默认值`1`