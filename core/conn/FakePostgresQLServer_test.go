package conn

import "testing"

func Test_start(t *testing.T) {
	host := "localhost"
	port := 10000
	targetOS := "windows"
	cmd := "calc.exe"
	path := "/xml"

	start(host, port, targetOS, cmd, path)
}
