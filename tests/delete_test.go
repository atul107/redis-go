package tests

import (
	"testing"
)

func TestDelHandler(t *testing.T) {
	t.Run("Delete Single Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "DEL mykey")
		assertResponse(t, response, "1")

		response = sendCommand(t, conn, "GET mykey")
		assertResponse(t, response, "(nil)")
	})

	t.Run("Delete Non-Existing Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "DEL nonexisting_key")
		assertResponse(t, response, "0")
	})

	t.Run("Delete Multiple Keys", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET key1 value1")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "SET key2 value2")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "DEL key1 key2")
		assertResponse(t, response, "2")

		response = sendCommand(t, conn, "GET key1")
		assertResponse(t, response, "(nil)")

		response = sendCommand(t, conn, "GET key2")
		assertResponse(t, response, "(nil)")
	})

	t.Log("All DELETE handler test cases passed.")
}
