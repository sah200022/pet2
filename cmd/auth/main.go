package main

import (
	"PetProject/internal/handler"
	"PetProject/internal/repository"
	"PetProject/internal/service"
	"fmt"
	"net/http"
)

func main() {

	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	fmt.Println("Запуск сервера")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("ошибка запуска сервера", err)
	}
}
