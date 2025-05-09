package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *Server) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqToService, err := req.ToService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.Service.InsertUserInSegments(r.Context(), reqToService)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	idUser := chi.URLParam(r, "id")
	id, err := uuid.Parse(idUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.UserID = id

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqToService, err := req.ToService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.Service.DeleteUserInSegments(r.Context(), reqToService)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleHistorySubscriptions(w http.ResponseWriter, r *http.Request) {
	var timeString HistorySubscription
	if err := json.NewDecoder(r.Body).Decode(&timeString); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	date, err := timeString.ToService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	history, err := s.Service.GetHistorySubscriptions(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type response struct {
		UserID uuid.UUID  `json:"user_id"`
		Title  string     `json:"title"`
		Create time.Time  `json:"create"`
		Delete *time.Time `json:"delete"`
	}
	resp := make([]response, len(history))
	for i, R := range history {
		resp[i].UserID = R.UserID
		resp[i].Title = R.TitleSegment
		resp[i].Create = R.CreatedAt
		resp[i].Delete = R.DeletedAt
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCheckTTLSubscriptions(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Service.CheckTTLSubscriptions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
