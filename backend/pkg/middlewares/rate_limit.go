package middlewares

import (
	"database/sql"
	"net"
	"net/http"
	"time"

	db "01social/pkg/db/sqlite"
	"01social/pkg/utilities"
)

func RateLimit(handler http.HandlerFunc, minInterval time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			utilities.WriteJSON(w, 500, "invalid Ip adress", nil)
			return
		}

		var lastRequest time.Time

		err = db.Database.QueryRow(
			"SELECT last_request FROM rate_limits WHERE ip = ? AND route = ?", ip, r.URL.Path,
		).Scan(&lastRequest)

		if err == sql.ErrNoRows {
			_, err = db.Database.Exec(
				"INSERT INTO rate_limits (ip, route, last_request) VALUES (?, ?, ?)",
				ip, r.URL.Path, time.Now(),
			)
			if err != nil {
				utilities.WriteJSON(w, 500, "Internal Server Error", nil)

				return
			}
			handler(w, r)
			return
		} else if err != nil {
			utilities.WriteJSON(w, 500, "Internal Server Error", nil)

			return
		}

		if time.Since(lastRequest) < minInterval {
			utilities.WriteJSON(w, 500, "Please wait before sending another request.", nil)

			return
		}

		_, err = db.Database.Exec(
			"UPDATE rate_limits SET last_request = ? WHERE ip = ? AND route = ?",
			time.Now(), ip, r.URL.Path,
		)
		if err != nil {
			utilities.WriteJSON(w, 500, "Internal Server Error", nil)

			return
		}

		handler(w, r)
	}
}
