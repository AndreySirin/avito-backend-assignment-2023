package server

import (
	"context"
	"encoding/json"
	"myapp/stor"
	"net/http"
)

func (s *Server) UserAddSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var subs stor.Subscription
	err := json.NewDecoder(r.Body).Decode(&subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.Sub.InsertUserInSegment(context.Background(), subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
