package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type requests struct {
	lastSenn time.Time
	count    int
}

var (
	visitors = make(map[string]*requests)
	muu      sync.Mutex
)

type rateLimitResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func writeRateLimitResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(rateLimitResponse{
		Message: message,
		Data:    data,
	})
}

func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		muu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &requests{time.Now(), 1}
			muu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSenn) < time.Minute {
			if v.count >= 20 {
				muu.Unlock()
				writeRateLimitResponse(w, http.StatusTooManyRequests, "Too manny requests, Tty after 1min", nil)
				return
			}
			v.count++
		} else {
			v.count = 1
			v.lastSenn = time.Now()
		}
		muu.Unlock()
		next.ServeHTTP(w, r)
	}
}
