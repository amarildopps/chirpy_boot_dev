package main

import "net/http"

func adminMetricsHandlerOld(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
