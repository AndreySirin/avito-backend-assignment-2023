package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/storage"
)

func (s *Server) UserAddSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var subs storage.Subscription
	err := json.NewDecoder(r.Body).Decode(&subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.Sub.InsertUserInSegment(context.Background(), subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode("successfully"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) UserDeleteSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var subs storage.Subscription
	err := json.NewDecoder(r.Body).Decode(&subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.Sub.DeleteUserInSegment(context.Background(), subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode("successfully"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
