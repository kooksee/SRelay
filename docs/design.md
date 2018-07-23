
# 项目设计

## 协议

### 添加客户端到srelay

```
{
	"event":"account",
	"account":"客户端名称或者ID",
	"token":"srelay的密码"
}
```
```
{
	"event":"account",
	"account":"123456",
	"token":"123456"
}
```

### 向客户端发送消息

```
{
	"event":"ws",
	"account":"客户端名称或者ID",
	"msg":"需要发送的消息"
}
```
```
{
	"event":"tcp",
	"account":"客户端名称或者ID",
	"msg":"需要发送的消息"
}
```

1. server 放到网络中
2. 客户端链接server端，并获取一个端口，组成一个外网的udp的addr
3. 客户端通过server 的一个ip获取