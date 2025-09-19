package handler

import (
	"encoding/json"
	"forum1/internal/entity"
	"forum1/internal/service"
	"html/template"
	"net/http"
	"strconv"
)

type ClubHandler struct {
	service service.ClubService
}

func NewClubHandler(s service.ClubService) *ClubHandler {
	return &ClubHandler{service: s}
}

// POST /clubs
func (h *ClubHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c entity.Club
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// GET /clubs/{id}
func (h *ClubHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // зависит от роутера, можно chi mux / gorilla mux
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// GET /clubs
func (h *ClubHandler) List(w http.ResponseWriter, r *http.Request) {
	clubs, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubs)
}

type ClubPageHandler struct {
	service   service.ClubService
	templates *template.Template
}

func NewClubPageHandler(s service.ClubService, tmpl *template.Template) *ClubPageHandler {
	return &ClubPageHandler{service: s, templates: tmpl}
}

// GET /boards/club
func (h *ClubPageHandler) ListPage(w http.ResponseWriter, r *http.Request) {
	clubs, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.templates.ExecuteTemplate(w, "clubs", clubs)
}

// GET /boards/club/{id}
func (h *ClubPageHandler) DetailPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // зависит от роутера
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	club, err := h.service.GetByID(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.templates.ExecuteTemplate(w, "club_detail", club)
}
