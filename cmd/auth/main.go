package main

import (
	ArticleHandler "PetProject/internal/articles/handler"
	ArticleRepository "PetProject/internal/articles/repository"
	ArticleService "PetProject/internal/articles/service"
	"PetProject/internal/config"
	"PetProject/internal/database"
	"PetProject/internal/middleware"
	"PetProject/internal/user/handler"
	"PetProject/internal/user/repository"
	"PetProject/internal/user/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	JWT_Middleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	cfg := config.Load()

	//Подключение к БД
	dbPool, err := database.NewPostgresPool(cfg.DB_URL)
	if err != nil {
		log.Fatal(err, "Failed to connect to DB")
	}
	defer dbPool.Close()

	r := chi.NewRouter()
	r.Use(JWT_Middleware.Logger)
	r.Use(JWT_Middleware.Recoverer)

	userRepo := repository.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo, []byte(cfg.JWT_SECRET))
	authHandler := handler.NewAuthHandler(authService)

	articleRepo := ArticleRepository.NewArticleRepository(dbPool)
	articleService := ArticleService.NewArticleService(articleRepo)
	articleHandler := ArticleHandler.NewArticleHandler(articleService)

	jwtMiddleware := middleware.JWTMiddleware([]byte(cfg.JWT_SECRET))

	//Auth
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.With(jwtMiddleware).Get("/me", authHandler.Me)
	})

	//Articles
	r.Route("/articles", func(r chi.Router) {

		r.Get("/", articleHandler.GetAll)
		r.Get("/{id}", articleHandler.GetID)

		r.With(jwtMiddleware).Post("/create", articleHandler.Create)
		r.With(jwtMiddleware).Delete("/{id}", articleHandler.Delete)
	})

	//Запуск сервера
	fmt.Println("Запуск сервера")
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		fmt.Println("ошибка запуска сервера", err)
	}
}
