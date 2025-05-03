package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]map[string]*visitor)
var mu sync.Mutex

var pathLimits = map[string]*rate.Limiter{
	"/api/login":     rate.NewLimiter(0.1, 3),
	"/api/register":  rate.NewLimiter(0.1, 3),
	"/api/cart":      rate.NewLimiter(1, 2),
	"/api/cart/bulk": rate.NewLimiter(0.5, 1),
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, pathMap := range visitors {
			for path, v := range pathMap {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(pathMap, path)
				}
			}
			if len(pathMap) == 0 {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func init() {
	go cleanupVisitors()
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return strings.TrimSpace(ip)
}

func getLimiter(ip, path string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := visitors[ip]; !ok {
		visitors[ip] = make(map[string]*visitor)
	}

	if v, ok := visitors[ip][path]; ok {
		v.lastSeen = time.Now()
		return v.limiter
	}

	limit, ok := pathLimits[path]
	if !ok {
		limit = rate.NewLimiter(1, 5)
	}

	visitors[ip][path] = &visitor{
		limiter:  rate.NewLimiter(limit.Limit(), limit.Burst()),
		lastSeen: time.Now(),
	}
	return visitors[ip][path].limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		path := r.URL.Path

		limiter := getLimiter(ip, path)

		if !limiter.Allow() {
			http.Error(w, "429 - Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
