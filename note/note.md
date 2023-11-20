## 构造服务端通信包

首先分析mysql初始发包到服务端，抓包分析

![mysqlpacket_1](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_1.png)

可以看到以下几步

1、greeting包，服务端返回了banner，其中包含mysql的版本

2、客户端登录请求

3、初始化查询

为了要实现一个fake mysql，需要自主构造greeting



**分析greeting**

```
4字节的包头 + 若干字节的包体。

解析包头字段：
 3字节：包体的有效长度，小端格式。
 1字节：包号。由于有效长度最大只能表示16M长度，所以增加一个包号字段，循环使用。

解析包体字段：
 1字节：协议版本号。
 若干字节：以0x00结尾的服务器版本描述。
 4字节：连接mysql使用的线程ID。
 9字节：8字节用于安全认证的随机数，+1字节的0x00填充。
 2字节：服务器能力标志集合。
 1字节：服务器默认使用的字符编码。
 2字节：服务器状态标志集合。
 13字节：保留。
 13字节：12字节用于安全认证的随机数补充，+1字节的0x00填充。
 若干字节：以0x00结尾的若干警告或其他描述文本。

关于安全认证的随机数：
 考虑到兼容性，随机数被分成2部分，旧版本使用8字节的随机数，新版本使用20字节的随机数。
 
看下面得结构更简洁
 
协议号 - Int<1> （Protocol number – Int<1>）
服务端版本号 - String （Server version – String）
线程 id - Int<4> （Thread id – Int<4>）
盐值1 - String （Salt1 – String）
服务端能力 - Int<2> （Server capabilities – Int<2>）
服务端语言 - Int<1> （Server language – Int<1>）
服务端状态 - Int<2> （Server Status – Int<2>）
扩展服务端能力 - Int<2> （Extended Server Capabilities – Int<2>）
身份认证插件长度 - Int<1> （Authentication plugin length – Int<1>）
保留字节 - 10 字节 （Reserved bytes – 10 bytes）
盐值2 - String （Salt2 – String）
身份验证插件字符串 - String （Authentication plugin string – String）
```

其中navicat的插件authplugin有区别，在返回相应包的时候需要返回切换插件的请求，也就是这点是通过用户输入请求是否包含“@clear”和客户端返回的插件同时判断的

![image-20231118093810618](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_2.png)



![image-20231118094016678](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_3.png)



而切换的数据包也十分好构造，返回eof数据包同时写入auth method name

客户端登录请求更简单，通过这个就明白了

![image-20231118094311413](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_4.png)

之后进行循环处理客户端请求，对于readfile来说，构造一个oxfb 返回code的packet和文件也是很简单的。



![image-20231118102236925](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_6.png)

可以看到local读取成功

![image-20231118102134738](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_5.png)



修改后能处理大文件了，只要不超过go的byte数组限制，也就是你运行的最大内存，这也是个很大的数字了，下面写实现读取大文件的思路

## 读取大文件

由于始终是读取mysql proto类型的包，需要写一个mysql packet的结构体，isEnd作为符号符，由于go中io.reader接口有con，直接设置mysql流类型为io.reaeder，每个数据包根据前面length判断，如果是个满的数据包，也就是



![image-20231119141807431](D:\Code Check\Go\ApricusFakemysql\note\mysqlpacket_7.png)



8192,说明包没读完，设置follow也就是跟随符号符是否继续下一个包的读取，之后循环读取包数据存在byte数组里。这里我尝试了100mb文件传输依旧没有任何问题。

