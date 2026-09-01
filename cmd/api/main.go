package main

import (
	"encoding/json"
	codes "github.com/stanleyHayes/ghanacodes/internal/codes"
	"log"
	"net/http"
	"os"
	"strings"
)

type server struct{ resolver *codes.Resolver }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://codes.digitalghana.dev")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func (s server) resolve(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/resolve/"), "/")
	if len(parts) != 2 {
		writeJSON(w, 400, map[string]string{"error": "path must contain namespace and code"})
		return
	}
	writeJSON(w, 200, s.resolver.Resolve(parts[0], parts[1], r.URL.Query().Get("at")))
}
func (s server) crosswalk(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	writeJSON(w, 200, s.resolver.Crosswalk(query.Get("fromNamespace"), query.Get("code"), query.Get("toNamespace"), query.Get("at")))
}
func (s server) bulk(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Queries []struct {
			Namespace string `json:"namespace"`
			Code      string `json:"code"`
			At        string `json:"at"`
		} `json:"queries"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&payload); err != nil || len(payload.Queries) > 1000 {
		writeJSON(w, 400, map[string]string{"error": "invalid request or more than 1000 queries"})
		return
	}
	results := make([]codes.Resolution, len(payload.Queries))
	for index, query := range payload.Queries {
		results[index] = s.resolver.Resolve(query.Namespace, query.Code, query.At)
	}
	writeJSON(w, 200, map[string]any{"version": s.resolver.Dataset().Version, "results": results})
}
func (s server) graphql(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10)).Decode(&payload); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid GraphQL request"})
		return
	}
	value := func(name string) string { result, _ := payload.Variables[name].(string); return result }
	if strings.Contains(payload.Query, "crosswalk") {
		result := s.resolver.Crosswalk(value("fromNamespace"), value("code"), value("toNamespace"), value("at"))
		writeJSON(w, 200, map[string]any{"data": map[string]any{"crosswalk": result}})
		return
	}
	if strings.Contains(payload.Query, "resolveCode") {
		result := s.resolver.Resolve(value("namespace"), value("code"), value("at"))
		writeJSON(w, 200, map[string]any{"data": map[string]any{"resolveCode": result}})
		return
	}
	writeJSON(w, 400, map[string]string{"error": "supported operations are resolveCode and crosswalk"})
}
func main() {
	resolver, err := codes.Load("data/codes.json")
	if err != nil {
		log.Fatal(err)
	}
	app := server{resolver}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]any{"status": "ok", "version": resolver.Dataset().Version, "codes": len(resolver.Dataset().Codes)})
	})
	mux.HandleFunc("GET /v1/namespaces", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]any{"version": resolver.Dataset().Version, "namespaces": resolver.Dataset().Namespaces})
	})
	mux.HandleFunc("GET /v1/resolve/", app.resolve)
	mux.HandleFunc("GET /v1/crosswalk", app.crosswalk)
	mux.HandleFunc("POST /v1/bulk-resolve", app.bulk)
	mux.HandleFunc("POST /graphql", app.graphql)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
