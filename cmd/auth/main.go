package main

import (
	handler2 "PetProject/internal/articles/handler"
	repository2 "PetProject/internal/articles/repository"
	service2 "PetProject/internal/articles/service"
	"PetProject/internal/database"
	"PetProject/internal/handler"
	"PetProject/internal/middleware"
	"PetProject/internal/repository"
	"PetProject/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {

	//Подключение к БД
	dbPool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal(err, "Failed to connect to DB")
	}
	defer dbPool.Close()

	mux := http.NewServeMux()
	userRepo := repository.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	articleRepo := repository2.NewArticleRepository(dbPool)
	articleService := service2.NewArticleService(articleRepo)
	articleHandler := handler2.NewArticleHandler(articleService)

	//Хэндлеры
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.Handle("/me", middleware.JWTMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/article/create", middleware.JWTMiddleware(http.HandlerFunc(articleHandler.Create)))
	mux.HandleFunc("/articles", articleHandler.GetAll)
	mux.HandleFunc("/articles/", articleHandler.GetID)
	mux.Handle("/delete/", middleware.JWTMiddleware(http.HandlerFunc(articleHandler.Delete)))

	//Запуск сервера
	fmt.Println("Запуск сервера")
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Println("ошибка запуска сервера", err)
	}
}
