package tests

import "testing"

//TestZrangeHandler ...
func TestZrangeHandler(t *testing.T) {
	t.Run("Simple ZRANGE by Score", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "ZADD myset 1 member1")
		sendCommand(t, conn, "ZADD myset 2 member2")
		sendCommand(t, conn, "ZADD myset -1 member3")

		response := sendCommand(t, conn, "ZRANGE myset 0 -1")

		assertResponse(t, response, "member3, member1, member2")
	})

	t.Log("All ZRANGE handler test cases passed.")
}
