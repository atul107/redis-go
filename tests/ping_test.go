package tests

import "testing"

func TestPingHandler(t *testing.T) {
	t.Run("PING with No Arguments", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "PING")

		expectedResponse := "PONG"
		assertResponse(t, response, expectedResponse)
	})

	t.Run("PING with Custom Message", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "PING HelloRedis")

		expectedResponse := "HelloRedis"
		assertResponse(t, response, expectedResponse)
	})

	t.Log("All PING handler test cases passed.")
}
