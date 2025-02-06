package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(addr string) error {
	http.Handle("/metrics", withCors(promhttp.Handler()))
	return http.ListenAndServe(addr, nil)
}

func withCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions { // preflight requests
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
