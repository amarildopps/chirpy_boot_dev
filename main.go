package main

import (
	"log"
	"net/http"
)

func main() {

	apiConfig := apiConfig{}
	filepathRoot, port := ".", "8080"
	fileServerAddr := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()

	// static files
	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(fileServerAddr))
	mux.HandleFunc("GET /admin/metrics", apiConfig.adminMetricsHandler)
	// api endpoints
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/metrics", apiConfig.countMetrics)
	mux.HandleFunc("/api/reset", apiConfig.resetMetrics)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
