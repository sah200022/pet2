package service

import (
	"PetProject/internal/articles/repository"
	"errors"
)

type ArticleService struct {
	articleRepo *repository.ArticleRepository
}

func NewArticleService(articleRepo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

func (a *ArticleService) Create(title string, text string, author string) error {
	if title == "" || text == "" {
		return errors.New("title or text is empty")
	}
	article := repository.Article{
		Title:  title,
		Text:   text,
		Author: author,
	}
	_, err := a.articleRepo.Create(&article)
	if err != nil {
		return err
	}
	return nil
}
