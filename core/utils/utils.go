package utils

import (
	"encoding/binary"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ToLittleEndian(input []byte) []byte {
	output := make([]byte, len(input))
	for i, j := 0, len(input)-1; i < len(input); i, j = i+1, j-1 {
		output[i] = input[j]
	}
	return output
}
func IntToBytes(value int) []byte {
	bytes := make([]byte, 3)
	bytes[0] = byte(value >> 16 & 0xFF)
	bytes[1] = byte(value >> 8 & 0xFF)
	bytes[2] = byte(value & 0xFF)
	return bytes
}
func StartsWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}
func CheckAndRemove(s, substr string) (bool, string) {
	index := strings.Index(s, substr)
	if index != -1 {
		clearRemovedString := strings.Replace(s, substr, "", 1)
		return true, clearRemovedString
	}
	return false, s
}
func ExtractFileName(filePath string) string {
	fileName := filepath.Base(filePath)
	return fileName
}
func ChangeInt(data int64) []byte {
	if data < 0 {
		panic("数据必须是非负数")
	} else if data < 251 {
		return []byte{byte(data)}
	} else if data < int64(1<<16) {
		buffer := make([]byte, 3)
		buffer[0] = byte(0xFC)
		binary.LittleEndian.PutUint16(buffer[1:], uint16(data))
		return buffer
	} else if data < int64(1<<24) {
		buffer := make([]byte, 4)
		buffer[0] = byte(0xFD)
		binary.LittleEndian.PutUint16(buffer[1:], uint16(data))
		buffer[3] = byte(data >> 16)
		return buffer
	} else if data < int64(1<<56) {
		buffer := make([]byte, 9)
		buffer[0] = byte(0xFE)
		binary.LittleEndian.PutUint64(buffer[1:], uint64(data))
		return buffer
	} else {
		panic("数据超过范围")
	}
}

func GetJavaBinPath() (string, error) {
	javaBinPath, err := exec.LookPath("java")
	if err != nil {
		return "", fmt.Errorf("无法找到Java可执行文件: %v", err)
	}

	return javaBinPath, nil
}
func GetYsoContent(ysoType string, command string) ([]byte, error) {
	javaBinPath, err := GetJavaBinPath()
	if err != nil {
		return nil, err
	}
	YsoGadgetPath := "../../plugin/YsoGadget.jar"

	cmd := exec.Command(javaBinPath, "-jar", YsoGadgetPath, ysoType, command)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}
func ExtractMySQLVersion(input string) (bool, string) {
	regex := regexp.MustCompile(`mysql-connector-java-(.+?)\s`)
	match := regex.FindStringSubmatch(input)
	if len(match) < 2 {
		fmt.Errorf("无法提取MySQL版本号")
		return false, ""
	}
	return regex.MatchString(input), match[1]
}
