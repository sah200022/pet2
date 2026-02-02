package repository

import "errors"

type Article struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

type ArticleRepository struct {
	articles map[int]*Article
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{
		articles: make(map[int]*Article),
	}
}

//Что умеет делать БД / Методы ArticleRepository

// Создание записи
func (a *ArticleRepository) Create(article *Article) (Article, error) {
	a.articles[article.ID] = article
	return *article, nil
}

// Найти статью по id
func (a *ArticleRepository) Get(id int) (Article, error) {
	article, ok := a.articles[id]
	if !ok {
		return Article{}, errors.New("article not found")
	}
	return *article, nil
}

// показать все статьи
func (a *ArticleRepository) GetAll() ([]Article, error) {
	var articles []Article
	for _, article := range a.articles {
		articles = append(articles, *article)
	}
	return articles, nil
}
