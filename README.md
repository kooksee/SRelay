# srelay

本版本暂时并没有重要的功能

> 超级中继

> 现在p2p项目很多无法完成，因为没有固定的外网IP，所以，通信很困难
> 那么我想，通过在一台主机上面通过端口转发的方式转发数据，虽然消耗了一点性能，但能够实现网络数据转发
> 该项目后端客户端有tcp和websocket两种方式的
> 前端的方式有http,udp,wensocket,tcp这几种模式的,以后可能会实现kcp的

## 项目初始化

### 请安装task项目管理工具

```shell
go get -u -v github.com/go-task/task/cmd/task
```

### 安装依赖

```shell
task deps

其他依赖请把依赖添加到scripts/deps.sh中
```



### Goland作为开发工具配置

    请把当前项目根目录设置为`Project GOPATH`

## 编译

```shell
task build
```

## 运行
```shell
task dev
```


检测能否打开外网IP port
如果不能够使用upnp pmp等方式暴露自己的端口，那么就通过中继的方式
中继服务端通过认证的方式连接
客户端通过用户名和密码登陆之后，
并通过接口向服务端发布自己的端口(网络模式,端口号),失败后重试


## docs

[协议设计](./docs/design.md)
