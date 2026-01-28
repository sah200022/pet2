package main

import (
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

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.Handle("/me", middleware.JWTMiddleware(http.HandlerFunc(authHandler.Me)))

	fmt.Println("Запуск сервера")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Println("ошибка запуска сервера", err)
	}
}
