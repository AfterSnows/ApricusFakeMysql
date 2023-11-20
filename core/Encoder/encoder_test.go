package Encoder

import (
	"ApricusFakemysql/core/utils"
	"fmt"
	"strings"
	"testing"
)

var Username string

func TestEncodeLoader(t *testing.T) {
	//result := EncodeLoader("Base64@md5@base64", "readFile=D://1.txt")
	//fmt.Println(result)

	SwitchClear, Usernamed := utils.CheckAndRemove("Base64@md5@base64:cmVhZEZpbGU9RDovLzEudHh0@clear", "@clear")
	fmt.Println(SwitchClear)
	IfEncode, EncodeString := utils.CheckAndRemove(Usernamed, "Encode=")
	fmt.Println(EncodeString)
	if IfEncode {
		//Encode=Base64@md5@base64:ZnVja3lvdQ==
		EncodeThing := strings.Split(EncodeString, ":")
		fmt.Println(len(EncodeThing))
		Username = DecodeLoader(EncodeThing[0], EncodeThing[1])
	} else {
		Username = Usernamed
	}
	fmt.Println(Username)
	IfRead, Filename := utils.CheckAndRemove(Username, "readFile=")
	fmt.Println(IfRead)
	fmt.Println(Filename)
}
