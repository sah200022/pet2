package handler

import (
	"PetProject/internal/articles/service"
	"PetProject/internal/middleware"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	articles, err := h.articleService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}
	json.NewEncoder(w).Encode(articles)
}

func (h *ArticleHandler) GetID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/articles/")

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
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, ok := r.Context().Value(middleware.ContextKeyUserID).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/delete/")
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
