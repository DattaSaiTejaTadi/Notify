package main

import (
	"Notify/config"
	"Notify/handler"
	"Notify/service"
	"Notify/store"
	"log"
	"net/http"
)

func main() {
	db := config.ConnectDB()
	userStore := &store.UserStore{DB: db}
	userService := &service.UserService{Store: userStore}
	userHandler := &handler.UserHandler{Service: userService}

	http.HandleFunc("/register", userHandler.RegisterUser)
	http.HandleFunc("/login", userHandler.LoginUser)
	http.HandleFunc("/protected", userHandler.AuthenticatedEndpoint)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
