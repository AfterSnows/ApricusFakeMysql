package Encoder

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"strings"
)

var EncoderLoaderMap = map[string]Encoder{
	"default": new(DefaultEncoder),
	"md5":     new(Md5Encoder),
	"base64":  new(Base64Encoder),
}

type Encoder interface {
	Encode(string) string
	Decode(string) string
}
type DefaultEncoder struct {
}

func (Encoder *DefaultEncoder) Encode(value string) string {
	return value
}
func (Encoder *DefaultEncoder) Decode(value string) string {
	return value
}

type Md5Encoder struct {
}

func (Encoder *Md5Encoder) Encode(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}

func (Encoder *Md5Encoder) Decode(value string) string {
	decoded, err := hex.DecodeString(value)
	if err != nil {
		// 处理解码错误
	}
	return string(decoded)
}

type Base64Encoder struct {
}

func (Encoder *Base64Encoder) Encode(value string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	return encoded
}

func (Encoder *Base64Encoder) Decode(value string) string {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		// 处理解码错误
	}
	return string(decoded)
}

func EncodeLoader(value string, data string) string {
	var EncodeResult string
	if strings.Contains(value, "@") {
		split := strings.Split(value, "@")
		splitLen := len(split)
		for i := 0; i < splitLen; i++ {
			if encode, ok := EncoderLoaderMap[strings.ToLower(split[i])]; ok {
				EncodeResult = encode.Encode(data)
			} else {
				log.Fatalf("Unknown 	Encodetype Error : %s\n", value)
			}
		}
	} else {
		if encode, ok := EncoderLoaderMap[strings.ToLower(value)]; ok {
			return encode.Encode(data)
		} else {
			log.Fatalf("Unknown 	Encodetype Error : %s\n", value)
		}
	}
	return EncodeResult
}
func DecodeLoader(value string, data string) string {
	var DecodeResult string
	if strings.Contains(value, "@") {
		split := strings.Split(value, "@")
		splitLen := len(split)
		for i := 0; i < splitLen; i++ {
			if encode, ok := EncoderLoaderMap[strings.ToLower(split[i])]; ok {
				DecodeResult = encode.Decode(data)
			} else {
				log.Fatalf("Unknown 	Encodetype Error : %s\n", value)
			}
		}
	} else {
		if encode, ok := EncoderLoaderMap[strings.ToLower(value)]; ok {
			return encode.Encode(data)
		} else {
			log.Fatalf("Unknown 	Encodetype Error : %s\n", value)
		}
	}
	return DecodeResult
}
