package proto

import (
	"bytes"
	"math/rand"
	"reflect"
)

var MysqlVersion = []string{
	"5.0.4",
}

type Greeting struct {
	Protocol         []byte
	Version          []byte //服务器版本描述
	VersionEnd       []byte //服务器版本描述0x00结尾
	ServerThreadID   []byte
	Salt             []byte
	Padding          []byte
	ServerCapL       []byte
	ServerLanguage   []byte
	ServerStatus     []byte
	ExtendServerCapL []byte
	AuthPluginLength []byte
	Unused           []byte
	Salt2            []byte //用于安全认证的随机数补充
	Padding2         []byte
	AuthPlugin       []byte
	End              []byte //以0x00结尾的若干警告或其他描述文本
}

// todo NewGreetingMessage() 服务器状态和state如果全都是固定的会不会对mysql连接产生影响，如何解决这个问题
func NewGreetingMessage() []byte {
	greeting := Greeting{
		Protocol:         []byte{0x0a},
		Version:          []byte(MysqlVersion[rand.Intn(1)]),
		VersionEnd:       []byte{0x00},
		ServerThreadID:   []byte{0x00, 0x00, 0x00, 0x00},
		Salt:             []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Padding:          []byte{0x00},       //随机盐后的填充
		ServerCapL:       []byte{0xff, 0xff}, //表示服务器所支持的功能和特性的集合
		ServerLanguage:   []byte{0x21},       //utf-8
		ServerStatus:     []byte{0x02, 0x00},
		ExtendServerCapL: []byte{0x08, 0x00},
		AuthPluginLength: []byte{0x14},
		Unused:           []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Salt2:            []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x08, 0x07, 0x06},
		Padding2:         []byte{0x00},
		AuthPlugin:       []byte("mysql_clear_password"),
		End:              []byte{0x00},
	}
	serverGreetingBytes := greetingToBytes(greeting)
	return serverGreetingBytes
}
func greetingToBytes(greeting Greeting) []byte {
	s := reflect.ValueOf(greeting)
	buf := new(bytes.Buffer)

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		buf.Write(field.Bytes())
	}

	return buf.Bytes()
}
