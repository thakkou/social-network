package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"01social/pkg/utilities"
	"01social/pkg/ws"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandlerWs(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		http.Error(w, `{"error":"missing ticket"}`, http.StatusUnauthorized)
		return
	}

	id, err := utilities.RedeemTicket(ticket)
	if err != nil {
		http.Error(w, `{"error":"invalid or expired ticket"}`, http.StatusUnauthorized)
		return
	}

	userId := strconv.Itoa(id)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := ws.StoreClient(userId, conn)
	go ws.HandleClient(client)
}

func GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	users := make([]string, 0)
	ws.Mu.RLock()
	for userID := range ws.Clients {
		users = append(users, userID)
	}
	ws.Mu.RUnlock()

	sort.Strings(users)
	utilities.WriteJSON(w, http.StatusOK, "online users fetched", users)
}

func TestBroadcast(w http.ResponseWriter, r *http.Request) {
	ws.BroadcastExcept("", "test_event", map[string]string{
		"message": "hello everyone 👋",
	})
	w.Write([]byte("sent"))
}
