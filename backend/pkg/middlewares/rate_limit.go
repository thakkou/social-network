package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"01social/pkg/utilities"
)

// ── In-memory rate limiter ──

type rateEntry struct {
	lastRequest time.Time
}

var (
	rateMu    sync.Mutex
	rateStore = make(map[string]*rateEntry)
)

// rateKey builds a "ip:route" key for the in-memory store.
func rateKey(ip, route string) string {
	return ip + ":" + route
}

func RateLimit(handler http.HandlerFunc, minInterval time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			utilities.WriteJSON(w, 500, "invalid Ip address", nil)
			return
		}

		key := rateKey(ip, r.URL.Path)

		rateMu.Lock()
		entry, exists := rateStore[key]
		if !exists {
			// First request – create entry
			rateStore[key] = &rateEntry{lastRequest: time.Now()}
			rateMu.Unlock()
			handler(w, r)
			return
		}

		if time.Since(entry.lastRequest) < minInterval {
			rateMu.Unlock()
			utilities.WriteJSON(w, 429, "Please wait before sending another request.", nil)
			return
		}

		entry.lastRequest = time.Now()
		rateMu.Unlock()

		handler(w, r)
	}
}
