# 网络工具箱

## 端口扫描

```shell
nts scan 192.168.1.0/24 22
Opened: 192.168.1.1:22
Opened: 192.168.1.3:22
Opened: 192.168.1.4:22
Opened: 192.168.1.11:22
Opened: 192.168.1.15:22
```

- `-T` 扫描最大线程，默认值`1`
- `-t` 扫描超时时间（毫秒），默认值`200`

## TCP 测试

```shell
nts ping www.baidu.com 443
Reply from www.baidu.com:443 time=13ms
Reply from www.baidu.com:443 time=9ms
Reply from www.baidu.com:443 time=20ms
Reply from www.baidu.com:443 time=12ms
Success=4, Error=0, Max=20ms, Min=9ms, Avg=13ms
```

- `-C` PING 的次数，默认值`4`
- `-i` 每次 PING 的间隔（毫秒），默认值`1000`
- `-P` 指定协议，默认值`tcp`，可选项 `tcp|http`
- `-t` PING 超时时间
- `--tls` 使用 `TLS` 建立连接

### 协议说明

- `tcp` 建立 TCP 连接后即认为成功
- `http` 建立 TCP 连接后，主动发送一个 HTTP `HEAD`请求，并读取 HTTP 服务器响应

## 带宽测试

### 服务端

```shell
nts sts
time="2025-03-06 22:11:11" level=info msg="tcp server listening on :8080"
time="2025-03-06 22:11:11" level=info msg="quic server listening on :8080"
time="2025-03-06 22:11:46" level=info msg="tcp download from 192.168.1.2:59367"
time="2025-03-06 22:11:56" level=info msg="download finished in 192.168.1.2:59367"
```

- `-s` 服务器监控地址，默认值`:8080`

### 客户端

```shell
nts stc -s 192.168.1.1:8080
time="2025-03-06 22:11:46" level=info msg="tcp download testing to 192.168.1.1:8080"
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

- `-s` 服务器地址，默认值`localhost:8080`
- `-m` 测试模式，默认值`download`，可选项`download|upload`
- `-t` 测试时间（秒），默认值`10`
- `-P` 测试的协议，默认值`tcp`，可选项`tcp|quic`
- `-T` 测试并发连接，默认值`1`

## 性能测试

### 服务端

```shell
nts bts 
time="2025-07-17 23:44:26" level=info msg="tcp server listening on :8080"
time="2025-07-17 23:44:43" level=info msg="tcp message: bd70a3c0-7e4e-456d-9eba-a6ccf9300fcf from 192.168.1.2:59735"
time="2025-07-17 23:44:44" level=info msg="tcp message: d7a2029f-6440-43fb-806a-e05b4c8198d8 from 192.168.1.2:59735"
time="2025-07-17 23:44:45" level=info msg="tcp message: 95625b03-0e1b-48f0-baf9-9ee0368b5b42 from 192.168.1.2:59735"
time="2025-07-17 23:44:46" level=info msg="tcp message: 126f2d44-9c88-4646-902a-f6034260bf12 from 192.168.1.2:59735"
time="2025-07-17 23:44:47" level=info msg="tcp message: 0f1d6de2-ed3c-4c98-b345-4e5fa8c0377a from 192.168.1.2:59735"
time="2025-07-17 23:44:48" level=info msg="tcp message: 0333e336-f006-4543-b41c-930a23794f78 from 192.168.1.2:59735"
time="2025-07-17 23:44:49" level=info msg="tcp message: 9e514627-a746-493f-b4b5-d70364730475 from 192.168.1.2:59735"
```

- `-s` 服务器监控地址，默认值`:8080`
- `-P` 测试协议，默认值`tcp`，可选值`tcp|udp|http|https|ws|quic`
- `-t` 读取消息超时，默认值`60s`

### 服务端

```shell
nts btc -s 192.168.1.1:8080
time="2025-07-17 23:44:43" level=info msg="tcp message: bd70a3c0-7e4e-456d-9eba-a6ccf9300fcffrom 192.168.1.1:8080 10ms"
time="2025-07-17 23:44:45" level=info msg="tcp message: d7a2029f-6440-43fb-806a-e05b4c8198d8from 192.168.1.1:8080 10ms"
time="2025-07-17 23:44:46" level=info msg="tcp message: 95625b03-0e1b-48f0-baf9-9ee0368b5b42from 192.168.1.1:8080 10ms"
time="2025-07-17 23:44:47" level=info msg="tcp message: 126f2d44-9c88-4646-902a-f6034260bf12from 192.168.1.1:8080 9ms"
time="2025-07-17 23:44:48" level=info msg="tcp message: 0f1d6de2-ed3c-4c98-b345-4e5fa8c0377afrom 192.168.1.1:8080 9ms"
time="2025-07-17 23:44:49" level=info msg="tcp message: 0333e336-f006-4543-b41c-930a23794f78from 192.168.1.1:8080 10ms"
time="2025-07-17 23:44:50" level=info msg="tcp message: 9e514627-a746-493f-b4b5-d70364730475from 192.168.1.1:8080 10ms"

```

- `-s` 服务器地址，默认值`localhost:8080`
- `-P` 测试协议，默认值`tcp`，可选值`tcp|udp|http|https|ws|quic`
- `-t` 读取消息超时，默认值`60s`
- `-i` 发送消息间隔，默认值`1000ms`
- `-T` 测试并发连接，默认值`1`
- `-m` 最大发送消息数，默认值`0`，`0`表示不限制