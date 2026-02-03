package handler

import (
	"PetProject/internal/articles/service"
	"encoding/json"
	"net/http"
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ArticleResponse{Message: "Article Created"})
}
