package middleware

import "net/http"

// CORS allows cross origin requests to the server. Note: By default it allows
// all origins so can be insecure.
// TODO: Make the origins configurable
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			w.Header().Add("Access-Control-Expose-Headers", "X-Praelatus-Token, Content-Type")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Add("Accepts", "application/json")

			if r.Method == "OPTIONS" {
				w.Write([]byte{})
				return
			}

			next.ServeHTTP(w, r)
		})
}
