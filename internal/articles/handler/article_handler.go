package handler

import (
	"PetProject/internal/articles/service"
	"PetProject/internal/middleware"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ArticleRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ArticleResponse struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var article ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Json Decode Error",
		})
		return
	}
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}
	err := h.articleService.Create(article.Title, article.Text, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ArticleResponse{Message: "Article Created"})
}

func (h *ArticleHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	w.Header().Set("Content-Type", "application/json")

	articles, total, err := h.articleService.GetAll(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}
	response := map[string]interface{}{
		"data":  articles,
		"page":  page,
		"limit": limit,
		"total": total,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ArticleHandler) GetID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID",
		})
		return
	}

	article, err := h.articleService.GetID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Not Found",
		})
		return
	}
	json.NewEncoder(w).Encode(article)
}

func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {

	_, ok := r.Context().Value(middleware.ContextKeyUserID).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID",
		})
		return
	}
	err = h.articleService.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Not Found",
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
