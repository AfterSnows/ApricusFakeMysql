package payload

import (
	"log"
	"strings"
)

var Payload []byte
var PayLoaderMap = map[string]payloadLoader{
	"custom": new(CustomPayloadLoader),
	"yso":    new(YSOPayloadLoader),
}

// custom_isaiutgys
func NewPayload(value string) []byte {
	values := strings.Split(value, "_")
	splitLen := len(values)
	if splitLen != 2 {
		log.Fatal("Wrong Payload Input,Check Input")
	}
	LoadType := strings.ToLower(values[0])
	SourceValue := values[1]
	if loader, ok := PayLoaderMap[LoadType]; ok {
		Payload = loader.Load(SourceValue)
	} else {
		log.Fatalf("Unknown 	Payloadertype Error : %s\n", LoadType)
	}
	return Payload
}
