gopush
======

第三方推送平台 GO 语言客户端。

目前仅支持 jpush (http://www.jpush.cn)

api 依据 jpush 官方文档( http://docs.jpush.cn/display/dev/Push+API+v2 )写成，并测试通过 



安装


```
go get github.com/wuyongzhi/gopush
```





示例
------

```
	m := NewRequest()
	m.AppKey("app_key")
	m.SendNo(1)
	m.ReceiverType(ReceiverTypeBoardcast)
	m.Platform("android")
	m.Sign("master_secret")
	m.MsgType(MsgTypeNotify)
	m.MsgContent(0, "", "hello,world", "")

	response, err := m.Send()

```


