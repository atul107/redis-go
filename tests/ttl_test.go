package tests

import (
	"testing"
	"time"
)

// TestTTLHandler ..
func TestTTLHandler(t *testing.T) {
	t.Run("TTL for Non-Existent Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "TTL non_existent_key")

		assertResponse(t, response, "-1")
	})

	t.Run("TTL for Key with No Expiry", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "SET noexpkey myvalue")

		response := sendCommand(t, conn, "TTL noexpkey")

		assertResponse(t, response, "-1")
	})

	t.Run("TTL for Expired Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "SET mykey myvalue EX 2")

		time.Sleep(3 * time.Second)

		response := sendCommand(t, conn, "TTL mykey")

		assertResponse(t, response, "-1")
	})

	t.Log("All TTL handler test cases passed.")
}
