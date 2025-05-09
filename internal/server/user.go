package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
)

func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Valid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqToService := req.ToService()

	id, err := s.Service.CreateUser(r.Context(), reqToService)
	if err != nil {
		if errors.Is(err, service.ErrNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	type response struct {
		ID uuid.UUID `json:"id"`
	}

	if err = json.NewEncoder(w).Encode(response{
		ID: id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	UserID := chi.URLParam(r, "id")
	id, err := uuid.Parse(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = id
	if err = req.Valid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqToService := req.ToService()

	err = s.Service.UpdateUser(r.Context(), reqToService)
	if err != nil {
		if errors.Is(err, service.ErrNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	type response struct {
		ID uuid.UUID `json:"id"`
	}
	if err = json.NewEncoder(w).Encode(response{
		ID: id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	UserID := chi.URLParam(r, "id")
	id, err := uuid.Parse(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = s.Service.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	type response struct {
		ID uuid.UUID `json:"id"`
	}
	if err = json.NewEncoder(w).Encode(response{
		ID: id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	UserID := chi.URLParam(r, "id")
	id, err := uuid.Parse(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	user, err := s.Service.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.Service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	UserID := chi.URLParam(r, "id")
	id, err := uuid.Parse(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subs, err := s.Service.GetUsersSubscription(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type response struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		IsAutoAdd bool      `json:"is_auto_add"`
	}

	resp := make([]response, len(subs))
	for i, _ := range subs {
		resp[i] = response{
			Title:     subs[i].TitleSegment,
			IsAutoAdd: subs[i].IsAutoAdded,
		}
	}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
