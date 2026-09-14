package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Server struct {
	config map[string]string
	mu     sync.RWMutex
}

type versionResponse struct {
	Version string `json:"version"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type environmentResponse struct {
	Environment string `json:"environment"`
}

type configResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type configRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type deleteResponse struct {
	Deleted bool `json:"deleted"`
}

func NewServer() *Server {

	return &Server{
		config: make(map[string]string),
	}
}

func (s *Server) Handler(environment string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /version", versionHandler)
	mux.HandleFunc("GET /env", environmentHandler(environment))
	mux.HandleFunc("POST /config", s.createConfigHandler)
	mux.HandleFunc("GET /config/{name}", s.getConfigHandler)
	mux.HandleFunc("DELETE /config/{name}", s.deleteConfigHandler)

	return mux
}

func writeJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	writeJson(w, http.StatusOK, versionResponse{
		Version: "1.0.0",
	})

}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJson(w, http.StatusOK, healthResponse{
		Status: "ok",
	})

}

func environmentHandler(environment string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJson(w, http.StatusOK, environmentResponse{
			Environment: environment,
		})
	}
}

func (s *Server) createConfigHandler(w http.ResponseWriter, r *http.Request) {
	var request configRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json request",
		})

		return
	}

	if request.Name == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.config[request.Name] = request.Value

	writeJson(w, http.StatusOK, configResponse{
		Name:  request.Name,
		Value: request.Value,
	})
}

func (s *Server) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	s.mu.RLock()
	defer s.mu.RUnlock()
	value, found := s.config[name]

	if !found {
		writeJson(w, http.StatusNotFound, map[string]string{
			"error": "value for this item not found",
		})
		return
	}
	writeJson(w, http.StatusOK, configResponse{
		Name:  name,
		Value: value,
	})

}

func (s *Server) deleteConfigHandler(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")

	s.mu.Lock()
	defer s.mu.Unlock()
	_, found := s.config[name]
	if found {
		delete(s.config, name)
	}

	if !found {
		writeJson(w, http.StatusNotFound, map[string]string{
			"error": "value for this item not found",
		})
		return
	}

	writeJson(w, http.StatusOK, deleteResponse{
		Deleted: true,
	})

}
