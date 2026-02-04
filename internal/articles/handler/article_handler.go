package handler

import (
	"PetProject/internal/articles/service"
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
		w.Write([]byte("Json Decode Error"))
		return
	}
	author, ok := r.Context().Value("email").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	err := h.articleService.Create(article.Title, article.Text, author)
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
		w.Write([]byte("Internal Server Error"))
		return
	}
	json.NewEncoder(w).Encode(articles)
}

func (h *ArticleHandler) GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/articles/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	article, err := h.articleService.GetID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)

}
