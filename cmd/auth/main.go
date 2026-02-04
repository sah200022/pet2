package main

import (
	handler2 "PetProject/internal/articles/handler"
	repository2 "PetProject/internal/articles/repository"
	service2 "PetProject/internal/articles/service"
	"PetProject/internal/handler"
	"PetProject/internal/middleware"
	"PetProject/internal/repository"
	"PetProject/internal/service"
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	articleRepo := repository2.NewArticleRepository()
	articleService := service2.NewArticleService(articleRepo)
	articleHandler := handler2.NewArticleHandler(articleService)

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.Handle("/me", middleware.JWTMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.HandleFunc("/article/create", articleHandler.Create)
	mux.HandleFunc("/articles", articleHandler.GetAll)
	mux.HandleFunc("/articles/", articleHandler.GetID)

	fmt.Println("Запуск сервера")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Println("ошибка запуска сервера", err)
	}
}
