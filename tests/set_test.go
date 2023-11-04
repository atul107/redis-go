package tests

import (
	"testing"
	"time"
)

// TestSetHandler ...
func TestSetHandler(t *testing.T) {
	// Start the server in a goroutine
	// go func() {
	// 	config := &Config{
	// 		TCP: TCP{
	// 			Addr: "127.0.0.1:6379",
	// 		},
	// 	}
	// 	Run(config)
	// }()
	// conn, err := net.Dial("tcp", "127.0.0.1:6379")
	// if err != nil {
	// 	t.Fatal("Error connecting to Redis server:", err)
	// }
	// defer conn.Close()

	// reader := NewRespReader(conn)
	// writer := bufio.NewWriter(conn)

	t.Run("Simple SET", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

	})

	t.Run("SET with NX option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykeyNX myvalue NX")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "SET mykeyNX myvalue NX")
		assertResponse(t, response, "(nil)")

	})

	t.Run("SET with XX option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykeyXX myvalue XX")
		assertResponse(t, response, "(nil)")

		response = sendCommand(t, conn, "SET mykeyXX myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "SET mykeyXX myvalue XX")
		assertResponse(t, response, "OK")

	})

	t.Run("SET with EX option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykeyEX myvalue EX 2")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "GET mykeyEX")
		assertResponse(t, response, "myvalue")

		time.Sleep(3 * time.Second)

		response = sendCommand(t, conn, "GET mykeyEX")
		assertResponse(t, response, "(nil)")
	})

	t.Run("SET with PX option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue PX 2000")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "GET mykey")
		assertResponse(t, response, "myvalue")

		time.Sleep(3 * time.Second)

		response = sendCommand(t, conn, "GET mykey")
		assertResponse(t, response, "(nil)")
	})

	t.Log("All SET handler test cases passed.")
}
