package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"redcoins/app"
	"redcoins/controllers"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/transactions/newsell", controllers.CreateTransactionSell).Methods("POST")
	router.HandleFunc("/api/transactions/newbuy", controllers.CreateTransactionBuy).Methods("POST")
	router.HandleFunc("/api/user/me/transactions", controllers.GetTransactionsForMe).Methods("GET")
	router.HandleFunc("/api/user/{userId}/transactions", controllers.GetTransactionsForUser).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":" + port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}