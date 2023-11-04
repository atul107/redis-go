package tests

import "testing"

//TestZaddHandler ...
func TestZaddHandler(t *testing.T) {
	t.Run("Simple ZADD", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "ZADD set 1.5 member1")

		assertResponse(t, response, "1")

		response = sendCommand(t, conn, "ZRANGE set 0 -1")

		assertResponse(t, response, "member1")
	})

	t.Run("ZADD with NX and CH Options", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "ZADD aset 1.5 member1")

		response := sendCommand(t, conn, "ZADD aset NX CH 2.5 member1")

		assertResponse(t, response, "0")

		response = sendCommand(t, conn, "ZRANGE aset 0 -1")

		assertResponse(t, response, "member1")
	})

	t.Log("All ZADD handler test cases passed.")
}
