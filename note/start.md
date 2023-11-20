MySQL服务端读取客户端文件漏洞利用


根据MySQL的官方文档，连接握手阶段中会执行如下操作：
客户端和服务端交换各自功能
如果需要则创建SSL通信通道
服务端认证客户端身份
身份认证通过后，客户端会在实际操作之前发送请求，等待服务器的响应。“Client Capabilities”报文中包括名为Can Use LOAD DATA LOCAL的一个条目
1、greeting包，服务端返回了banner，其中包含mysql的版本
2、客户端登录请求
3、然后是初始化查询，这里因为是phpmyadmin所以初始化查询比较多
4、load file local


note:
构建恶意服务端
1.回复mysql client一个greeting包 
2.等待client端发送一个查询包 
3.回复一个file transfer包
构造greeting包

许少写的好像有点问题，之后提一下建议