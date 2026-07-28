package ws

import (
	"sync"
	"time"
)

func HandleClient(client *Client) {
	// ── Periodic session validation (goroutine + mutex, no channels) ──
	var stopMu sync.Mutex
	stopped := false

	go func() {
		for {
			stopMu.Lock()
			if stopped {
				stopMu.Unlock()
				return
			}
			stopMu.Unlock()

			time.Sleep(30 * time.Second)

			stopMu.Lock()
			if stopped {
				stopMu.Unlock()
				return
			}
			stopMu.Unlock()

			if ValidateSession != nil && !ValidateSession(client.id) {
				// Session no longer valid – send force_logout and disconnect
				client.writeMu.Lock()
				client.conn.WriteJSON(map[string]any{
					"event_type": "force_logout",
					"data": map[string]string{
						"reason": "session expired",
					},
				})
				client.writeMu.Unlock()
				client.conn.Close()
				return
			}
		}
	}()

	defer func() {
		stopMu.Lock()
		stopped = true
		stopMu.Unlock()

		RemoveClient(client.id, client)
		client.conn.Close()

		Mu.RLock()
		_, stillOnline := Clients[client.id]
		Mu.RUnlock()

		if !stillOnline {
			BroadcastExcept(client.id, "client_disconnect", client.id)
		}
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		HandleMessage(client, msg)
	}
}
