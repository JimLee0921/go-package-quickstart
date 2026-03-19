## etcd

etcd 本质上就是一个分布式的 Key-Value 存储，类似于加强版的 redis，主要用于服务于微服务世界，相当于通讯录+配置中心+协调器。

### 核心作用


#### 服务注册于服务发现

go-zero 进行微服务架构创建时的服务注册于服务发现。

比如有多个服务：

```
user-service: 192.168.7.100:8080
order-service: 192.168.7.101:8080
```

这时 etcd 主要解决 user-service 与 order-service 进行通信时的问题, etcd 里面设置：

```
key: /services/user
value: 10.0.0.1:8080
```

服务启动时：

- user-service 注册到 etcd
- order-service 从 etcd 中查询

#### 配置中心

比如数据库等也可以在 etcd 中进行集中配置，所有服务从 etcd 中读取配置，而不是写死在代码中

```
/db/mysql = root:123456@tcp(...)
/redis/addr = 127.0.0.1:6379
```

### 下载

Windows 上直接在 GitHub 的 release 里面找到 Windows 安装包下载后解压即可

```
etcd.exe // 服务端
etcdctl.exe // 客户端
```

### 常见命令

```
etcdctl put key value              // 设置或修改键值对, 微服务通常使用 rpc.user, rpc.order 进行设置
ectdctl get key                    // 获取键值对(也可以判断 key 是否存在)
etcdctl get key --print-value-only // 只要值
etcdctl del key                    // 删除键值对
etcdctl get key --prefix           // 微服务最常用, 查询某个目录下所有key
etcdctl watch key                  // 监听键值对变化
etcdctl get "" --prefix            // 查看所有的 key
```