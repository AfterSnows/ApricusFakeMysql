package conn

import (
	"bytes"
	"log"
	"net"
	"strings"
)

type Client struct {
	Username string
}
type ClientResp struct {
	Command   int
	Statement string
}

func NewClientResp(RespData []byte) *ClientResp {
	if len(RespData) < 5 {
		log.Fatal("request too short")
		return nil
	}
	return &ClientResp{
		Command:   int(RespData[4]),
		Statement: string(RespData[5:]),
	}

}
func GetUsername(FirstResp []byte) string {
	var bao bytes.Buffer
	if len(FirstResp) < 40 {
		return ""
	}
	for i := 36; i < len(FirstResp); i++ {
		if FirstResp[i] == 0 {
			break
		}
		bao.Write([]byte{FirstResp[i]})
	}
	return bao.String()
}

func GetAuthPlugin(FirstResp []byte) string {
	var bao bytes.Buffer
	if len(FirstResp) < 40 {
		return ""
	}
	for i := 36; i < len(FirstResp); i++ {
		if FirstResp[i] == 0 {
			for i2 := i + 1; i2 < len(FirstResp); i2++ {
				if FirstResp[i2] == 0 {
					break
					break
				}
				bao.Write([]byte{FirstResp[i2]})
			}
		}

	}
	return bao.String()
}
func GetClientData(con net.Conn) []byte {
	RespBuffer := make([]byte, 1024*10)
	n, _ := con.Read(RespBuffer)
	ClientPostFirstData := RespBuffer[:n]
	return ClientPostFirstData
}

func PostClientData(conn net.Conn, PostData []byte) {
	_, err := conn.Write(PostData)
	if err != nil {
		log.Println(err)
		return
	}

}

func CheckAndRemove(s, substr string) (bool, string) {
	index := strings.Index(s, substr)
	if index != -1 {
		clearRemovedString := strings.Replace(s, substr, "", 1)
		return true, clearRemovedString
	}
	return false, s
}
