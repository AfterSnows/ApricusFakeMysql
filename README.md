# ApricusFakemysql

这是一个完全用go实现的fakemysql工具



## 利用

mysql local infile 读取

mysql jdbc 反序列化漏洞 RCE

mysql postgresql 反序列化漏洞 RCE



## 特点

1.支持文件名加密

2.支持大文件传输[取决于你的运行内存，应该能轻松1个g]

3.支持navicat[不管什么版本]

4.payload支持





## 点子

1.对于greeting的发包的mysql版本会随机发送版本

2.文件名支持加密防止特殊符号

## 运行测试

![image-20231119143828434](D:\Code Check\Go\ApricusFakemysql\asset\mysqlpacket_1.png)

利用成功【测试环境 java1.8+mysql-connector-java 8/5.1.27】 

![image-20231120162402833](D:\Code Check\Go\ApricusFakemysql\asset\mysqlpacket_2.png)







## note笔记

记录ApricusFakemysql中开发中学习的内容