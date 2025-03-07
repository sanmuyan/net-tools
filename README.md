# 网络工具箱

## 端口扫描

```shell
net-tools scan 192.168.1.0/24 22
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
net-tools ping www.baidu.com 443
Reply from www.baidu.com:443 time=13ms
Reply from www.baidu.com:443 time=9ms
Reply from www.baidu.com:443 time=20ms
Reply from www.baidu.com:443 time=12ms
Success=4, Error=0, Max=20ms, Min=9ms, Avg=13ms
```

- `-c` PING 的次数，默认值`4`
- `-i` 每次 PING 的间隔（毫秒），默认值`1000`
- `-P` 指定协议，默认值`tcp`，可选项 `tcp|http|https|read`
- `-T` PING 超时时间

### 协议说明

- `tcp` 建立TCP 连接后即认为成功
- `http` 主动发送一个 HTTP `HEAD`请求，并读取 HTTP 服务器响应
- `read` 建立连接后，读取一次服务端发送的`hello`包

## 带宽测试

### 服务端

```shell
net-tools sts
time="2025-03-06 22:11:11" level=info msg="tcp server listening on 0.0.0.0:8080"
time="2025-03-06 22:11:11" level=info msg="udp server listening on 0.0.0.0:8080"
time="2025-03-06 22:11:46" level=info msg="tcp download from 127.0.0.1:59367"
time="2025-03-06 22:11:56" level=info msg="download finished in 127.0.0.1:59367"
```

- `-s` 服务器监控地址，默认值`0.0.0.0:8080`
- `-P` 服务器监听协议，默认值`tcp-udp`，可选项`tcp-udp|tcp|udp`

### 服务端

```shell
net-tools stc 192.168.1.1:8080
time="2025-03-06 22:11:46" level=info msg="tcp download testing to localhost:8080"
time="2025-03-06 22:11:47" level=info msg="real-time speed: 339.38Mbps/s"
time="2025-03-06 22:11:48" level=info msg="real-time speed: 349.91Mbps/s"
time="2025-03-06 22:11:49" level=info msg="real-time speed: 344.06Mbps/s"
time="2025-03-06 22:11:50" level=info msg="real-time speed: 335.20Mbps/s"
time="2025-03-06 22:11:51" level=info msg="real-time speed: 354.65Mbps/s"
time="2025-03-06 22:11:52" level=info msg="real-time speed: 372.55Mbps/s"
time="2025-03-06 22:11:53" level=info msg="real-time speed: 351.18Mbps/s"
time="2025-03-06 22:11:54" level=info msg="real-time speed: 312.53Mbps/s"
time="2025-03-06 22:11:55" level=info msg="real-time speed: 299.41Mbps/s"
time="2025-03-06 22:11:56" level=info msg="real-time speed: 300.38Mbps/s"
time="2025-03-06 22:11:56" level=info msg="finished avg speed: 337.53Mbps/s"
```

- `-m` 测试模式，默认值`download`，可选项`download|upload`
- `-t` 测试时间（秒），默认值`10`
- `-P` 测试的协议，默认值`tcp`，可选项`tcp|udp`
- `-T` 测试并发连接，默认值`1`