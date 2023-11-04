package tests

import "testing"

func TestGetHandler(t *testing.T) {
	t.Run("Get Value of Existing Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "GET mykey")
		assertResponse(t, response, "myvalue")
	})

	t.Run("Get Value of Non-Existing Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "GET nonexisting_key")
		assertResponse(t, response, "(nil)")
	})

	t.Run("Get Value with Multiple Arguments", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "GET mykey extra_arg")
		assertResponse(t, response, "ERR wrong number of arguments for 'get' command")
	})
	t.Log("All GET handler test cases passed.")
}
