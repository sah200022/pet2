package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Article struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	AuthorID int    `json:"author_id"`
}

type ArticleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{
		db: db,
	}
}

//Что умеет делать БД / Методы ArticleRepository

// Создание записи
func (a *ArticleRepository) Create(article *Article) (Article, error) {
	query := `
	INSERT INTO articles (title, text, author_id)
	VALUES ($1, $2, $3)
	RETURNING id;
`
	var articleId int
	err := a.db.QueryRow(context.Background(), query, article.Title, article.Text, article.AuthorID).Scan(&articleId)
	if err != nil {
		return Article{}, err
	}
	article.ID = articleId
	return *article, nil
}

// Найти статью по id
func (a *ArticleRepository) GetID(id int) (Article, error) {
	query := `
	SELECT id, title, text, author_id
	FROM articles
	WHERE id = $1;
`
	var article Article
	err := a.db.QueryRow(context.Background(), query, id).Scan(
		&article.ID, &article.Title, &article.Text, &article.AuthorID,
	)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

// показать все статьи
func (a *ArticleRepository) GetAll() ([]Article, error) {
	query := `
	SELECT id, title, text, author_id
	FROM articles;
`
	rows, err := a.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(
			&article.ID, &article.Title, &article.Text, &article.AuthorID,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil

}
