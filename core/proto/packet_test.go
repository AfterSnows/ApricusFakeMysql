package proto

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestNewPacket(t *testing.T) {
	result := NewOptionByte("d")
	fmt.Println(hex.EncodeToString(result))
}
