package tests

import (
	"bufio"
	"net"
	"strings"
	"testing"

	. "github.com/redis-go/pkg/resp"
)

func TestRedisServer(t *testing.T) {
	// config := &Config{
	// 	TCP: TCP{
	// 		Addr: "localhost:6379",
	// 	},
	// }

	// go func() {
	// 	Run(config)
	// }()

	// time.Sleep(1 * time.Second)

	TestZaddHandler(t)
	TestZrangeHandler(t)
	TestGetHandler(t)
	TestHsetHandler(t)
	TestHgetHandler(t)
	TestHgetallHandler(t)
	TestKeysHandler(t)
	TestPingHandler(t)
	TestTTLHandler(t)
	TestExpireHandler(t)
	TestDelHandler(t)
}

func connectToServer(t *testing.T) (net.Conn, func()) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatal("Error connecting to Redis server:", err)
	}

	cleanup := func() {
		conn.Close()
	}

	return conn, cleanup
}

func sendCommand(t *testing.T, conn net.Conn, command string) string {
	reader := NewRespReader(conn)
	writer := bufio.NewWriter(conn)
	reader.SendCommand(writer, strings.Split(command, " ")...)
	response, err := reader.ReadResponse()
	if err != nil {
		t.Fatal("Error reading response:", err)
	}

	return response
}

func assertResponse(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected response: %s, Got response: %s", expected, actual)
	}
}
