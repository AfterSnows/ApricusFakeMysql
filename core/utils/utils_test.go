package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestStartsWith(t *testing.T) {
	content, err := GetYsoContent("Jdk7u21", "calc")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(content)
	}
}
