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
func (a *ArticleRepository) GetAll(limit, offset int) ([]Article, int, error) {
	var articles []Article
	var total int

	ArticleQuery := `
	SELECT id, title, text, author_id
	FROM articles
	ORDER BY id DESC LIMIT $1 OFFSET $2;
`
	countQuery := `
SELECT COUNT(*)
FROM articles
`
	err := a.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := a.db.Query(context.Background(), ArticleQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var article Article
		err := rows.Scan(
			&article.ID, &article.Title, &article.Text, &article.AuthorID,
		)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// удалить статью по id
func (a *ArticleRepository) Delete(id int) error {
	query := `
	DELETE FROM articles WHERE id = $1;
`
	_, err := a.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
