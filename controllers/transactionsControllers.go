package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"redcoins/models"
	u "redcoins/utils"
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

var CreateTransaction = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user") . (uint) //Grab the id of the user that send the request
	transaction := &models.Transaction{}

	err := json.NewDecoder(r.Body).Decode(transaction)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	//fmt.Println(transaction.Bitcoins)

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

	fmt.Println(responseObject.Data.Quotes.BRL.Price)

	transaction.User_Id_1 = user
	transaction.Convert_Rt = responseObject.Data.Quotes.BRL.Price
	transaction.Final_Value = responseObject.Data.Quotes.BRL.Price * transaction.Bitcoins
	resp := transaction.Create()
	u.Respond(w, resp)
}

var GetTransactionsFor = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user") . (uint)
	data := models.GetTransactions(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}