package ws

import (
	"time"
)

func HandleClient(client *Client) {
	// ── Periodic session validation ──
	stopTicker := make(chan struct{})
	defer close(stopTicker)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
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
			case <-stopTicker:
				return
			}
		}
	}()

	defer func() {
		// 1. Use your existing RemoveClient function instead of deleting the whole user map
		RemoveClient(client.id, client)
		client.conn.Close()

		// 2. Only broadcast disconnect if the user has no more active tabs open
		Mu.RLock()
		_, stillOnline := Clients[client.id]
		Mu.RUnlock()

		// fmt.Println("handling client")
		if !stillOnline {
			BroadcastExcept(client.id, "client_disconnect", client.id)
		}
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			// fmt.Println("client disconnected:", client.id)
			return
		}

		HandleMessage(client, msg)
	}
}
