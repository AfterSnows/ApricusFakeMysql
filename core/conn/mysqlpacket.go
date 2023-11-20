package conn

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type MysqlStreamSequence struct {
	seq byte
}

func (s *MysqlStreamSequence) getSeq() byte {
	return s.seq
}

func (s *MysqlStreamSequence) check(seq byte) {
	if seq != s.seq {
		panic("Wrong sequence")
	}
	fmt.Print(seq)
	s.incr()
}

func (s *MysqlStreamSequence) incr() {
	s.seq = (s.seq + 1) & 0xff
}

func (s *MysqlStreamSequence) reset() {
	s.seq = 0
}

func CreateMysqlStreamSequence() *MysqlStreamSequence {
	return &MysqlStreamSequence{seq: 2}
}

type MysqlPacketReader struct {
	stream io.Reader
	length int
	follow bool
}

func CreateMysqlPacketReader(stream io.Reader, seq *MysqlStreamSequence) *MysqlPacketReader {
	reader := bufio.NewReader(stream)
	return NewMysqlPacketReader(reader, seq)
}
func NewMysqlPacketReader(stream io.Reader, seq *MysqlStreamSequence) *MysqlPacketReader {
	return &MysqlPacketReader{
		stream: stream,
		length: 0,
		follow: true,
	}
}

func (r *MysqlPacketReader) readLeadData() []byte {
	ldata := make([]byte, 4)
	_, err := io.ReadFull(r.stream, ldata)
	if err != nil {
		log.Print("null")
	}

	l := int(ldata[0]) + (int(ldata[1]) << 8) + (int(ldata[2]) << 16)
	r.length = l
	if l < 8192 {
		r.follow = false
	}

	return ldata
}

func (r *MysqlPacketReader) Read() ([]byte, error) {
	if r.follow != true {
		return []byte{0x00}, nil
	}
	// 读取包头
	ldata := r.readLeadData()
	if ldata[0] == 0 && ldata[1] == 0 {
		log.Printf("No File Content,This could be a client problem requiring multiple retries")
		return nil, nil
	}
	if ldata == nil {
		return nil, io.EOF
	}

	// 数据包长度小于等于包头字节长度的情况
	if r.length <= len(ldata) {
		return nil, nil
	}

	// 计算真实的数据包长度（减去包头长度）
	dataSize := r.length
	data := make([]byte, dataSize)
	_, err := io.ReadFull(r.stream, data)
	if err != nil {
		return nil, err
	}

	r.length = 0 // 数据包已完整读取，重置长度

	return data, nil
}
