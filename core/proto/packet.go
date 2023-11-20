package proto

import (
	"ApricusFakemysql/core/utils"
	"bytes"
	"log"
)

func NewPacket(number int, content []byte) []byte {
	if len(content) > 16777216 {
		log.Fatal("packet is too long")
		return nil
	}
	buf := new(bytes.Buffer)
	packetLen := len(content)
	buf.Write(utils.ToLittleEndian(utils.IntToBytes(packetLen)))
	buf.Write([]byte{byte(number)})
	buf.Write(content)
	return buf.Bytes()
}

func NewOptionByte(column string) []byte {
	Result := bytes.Join([][]byte{
		utils.ChangeInt(int64(len("def"))),
		[]byte("def"),
		[]byte{0x00},
		utils.ChangeInt(int64(len("a"))),
		[]byte("a"),
		utils.ChangeInt(int64(len("a"))),
		[]byte("a"),
		utils.ChangeInt(int64(len(column))),
		[]byte(column),
		utils.ChangeInt(int64(len(column))),
		[]byte(column),
		[]byte{0x0c, 0x3f, 0x00, 0x1c, 0x00, 0x00, 0x00, 0xfc, 0xff, 0xff, 0x00, 0x00, 0x00}},
		[]byte(""))
	return Result
}
func NewOptionBytes(values [][]byte) []byte {
	var finalValues []byte
	for _, value := range values {
		LValue := bytes.Join([][]byte{
			utils.ChangeInt(int64(len(value))),
			value,
		}, []byte(""))
		finalValues = append(finalValues, LValue...)
	}
	return finalValues
}
