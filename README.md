# 网络工具箱

## 端口扫描

```shell
./portscan -i 192.168.1.0/24 -p 22
192.168.1.1:22
192.168.1.3:22
192.168.1.4:22
192.168.1.11:22
192.168.1.15:22
```

- `-p` 端口或端口范围，默认值`0-65535`
- `-i` IP地址或网段，默认值`127.0.0.1`
- `-t` 扫描最大线程，默认值`1`
- `-T` 扫描超时时间（毫秒），默认值`200`

## TCP 测试

```shell
./tcpping -h www.baidu.com -p 443
www.baidu.com:443 time=21ms
www.baidu.com:443 time=19ms
www.baidu.com:443 time=16ms
www.baidu.com:443 time=10ms
success=4, error=0, max=21ms, min=10ms, avg=16ms
```


- `-h` 目的主机，默认值`127.0.0.1`
- `-p` 目的端口，默认值`22`
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
./speedtests
2023/08/01 18:02:25 tcp server runing 0.0.0.0:8080
2023/08/01 18:03:02 tcp download from 192.168.1.2:50916
2023/08/01 18:03:12 tcp download finished in 192.168.1.2:50916
```

- `-p` 服务器监听端口，默认值`8080`
- `-s` 服务器监控地址，默认值`0.0.0.0`
- `-P` 服务器监听协议，默认值`tcp`，可选项`tcp|udp`

### 服务端

```shell
./speedtestc -s 192.168.1.1 -p 8080
2023/08/01 18:03:03 tcp download testing to 192.168.1.1:8080
2023/08/01 18:03:13 finished speed: 897Mbps/s
```

- `-p` 服务器端口，默认值`8080`
- `-s` 服务器地址，默认值`localhost`
- `-m` 测试模式，默认值`download`，可选项`download|upload`
- `-t` 测试时间（秒），默认值`10`
- `-P` 测试的协议，默认值`tcp`，可选项`tcp|udp`
- `-T` 测试并发连接，默认值`1`