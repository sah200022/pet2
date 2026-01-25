package main

import (
	"PetProject/internal/handler"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/register", handler.Register)
	//http.HandleFunc("/login")

	fmt.Println("Запуск сервера")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("ошибка запуска сервера", err)
	}
}
