package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"redcoins/models"
	u "redcoins/utils"
	"strconv"
)

type API_Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Quotes Quotes `json:"quotes"`
}

type Quotes struct {
	BRL BRL `json:"BRL"`
}

type BRL struct {
	Price float64 `json:"price"`
}

var aux = 0.0

var CreateTransactionSell func(w http.ResponseWriter, r *http.Request) = func(w http.ResponseWriter, r *http.Request) {
	user1 := r.Context().Value("user") . (uint) //Grab the id of the user that send the request

	transaction := &models.Transaction{}

	err := json.NewDecoder(r.Body).Decode(transaction)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	response, err := http.Get("https://api.coinmarketcap.com/v2/ticker/1/?convert=BRL")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject API_Response
	json.Unmarshal(responseData, &responseObject)

	//fmt.Println(responseObject.Data.Quotes.BRL.Price)

	transaction.User_Id_1 = user1

	transaction.Convert_Rt = responseObject.Data.Quotes.BRL.Price
	transaction.Final_Value = responseObject.Data.Quotes.BRL.Price * transaction.Bitcoins

	resp := transaction.CreateSell()
	u.Respond(w, resp)
}

var CreateTransactionBuy func(w http.ResponseWriter, r *http.Request) = func(w http.ResponseWriter, r *http.Request) {
	user1 := r.Context().Value("user") . (uint) //Grab the id of the user that send the request

	transaction := &models.Transaction{}

	err := json.NewDecoder(r.Body).Decode(transaction)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	response, err := http.Get("https://api.coinmarketcap.com/v2/ticker/1/?convert=BRL")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject API_Response
	json.Unmarshal(responseData, &responseObject)

	//fmt.Println(responseObject.Data.Quotes.BRL.Price)

	transaction.User_Id_1 = user1

	transaction.Convert_Rt = responseObject.Data.Quotes.BRL.Price
	transaction.Final_Value = responseObject.Data.Quotes.BRL.Price * transaction.Bitcoins

	resp := transaction.CreateBuy()
	u.Respond(w, resp)
}

var GetTransactionsForMe = func(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user") . (uint)
	fmt.Println(user_id)
	data := models.GetTransactions(user_id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetTransactionsForUser = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["userId"])
	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	data := models.GetTransactions(uint(id))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}