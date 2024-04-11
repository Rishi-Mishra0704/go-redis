package server

import "testing"

func TestParseCommand(t *testing.T) {
	msg := "3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	parseCommand(msg)
}
