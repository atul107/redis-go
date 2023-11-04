package tests

import "testing"

func TestKeysHandler(t *testing.T) {
	t.Run("KEYS with Matching Keys", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "SET key1 value1")
		sendCommand(t, conn, "SET key2 value2")
		sendCommand(t, conn, "SET mykey value3")

		response := sendCommand(t, conn, "KEYS key*")

		expectedResponse := "key1, key2"
		assertResponse(t, response, expectedResponse)
	})

	t.Run("KEYS with No Matching Keys", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "KEYS nonexistent*")

		expectedResponse := ""
		assertResponse(t, response, expectedResponse)
	})

	t.Run("KEYS with Incorrect Number of Arguments", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "KEYS key1 key2")
		assertResponse(t, response, "ERR wrong number of arguments for 'keys' command")
	})

	t.Log("All KEYS handler test cases passed.")
}
