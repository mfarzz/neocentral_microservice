package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewReverseProxy creates a reverse proxy that routes requests to the target URL.
func NewReverseProxy(targetRawURL string) http.Handler {
	targetURL, err := url.Parse(targetRawURL)
	if err != nil {
		log.Fatalf("❌ Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the response to handle errors gracefully
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("⚠️  Proxy error to %s: %v", targetURL.String(), err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"success": false, "message": "Service unavailable"}`))
	}

	return proxy
}
