package main

import (
	"Notify/handler"
	"Notify/middelware"
	"Notify/models"
	"Notify/service"
	"Notify/store"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal("Error while loading env")
	}
	dbUserName := os.Getenv("DB_Username")
	dbPassword := os.Getenv("DB_Password")
	store := store.New(dbUserName, dbPassword)
	err = store.ResetUserActiveStates()
	if err != nil {
		log.Fatal("Issue setting up the DB")
	}
	service := service.New(store)
	handler := handler.New(service)

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/home", middelware.ValidateJWT(handler.AdminHome, &models.Validations{Role: "admin"}))
	http.HandleFunc("/seniorHome", middelware.ValidateJWT(handler.L1Homepage, &models.Validations{Role: "senior"}))
	http.HandleFunc("/memberHome", middelware.ValidateJWT(handler.L2Homepage, &models.Validations{Role: "member"}))
	http.HandleFunc("/register", middelware.ValidateJWT(handler.RegisterUser, &models.Validations{Role: "admin"}))
	http.HandleFunc("/login", handler.LoginUser)
	http.HandleFunc("/protected", middelware.ValidateJWT(handler.AuthenticatedEndpoint, nil))
	http.HandleFunc("/getMembers", middelware.ValidateJWT(handler.GetMembers, &models.Validations{Role: "senior"}))
	http.HandleFunc("/logout", middelware.ValidateJWT(handler.Logout, nil))
	http.HandleFunc("/getNotification", middelware.ValidateJWT(handler.GetNotification, &models.Validations{Role: "member"}))
	http.HandleFunc("/notifyMember", middelware.ValidateJWT(handler.Notify, &models.Validations{Role: "senior"}))
	http.HandleFunc("/getAllUsers", middelware.ValidateJWT(handler.GetAllUsers, &models.Validations{Role: "admin"}))
	http.HandleFunc("/deleteUser", middelware.ValidateJWT(handler.DeleteUser, &models.Validations{Role: "admin"}))

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
