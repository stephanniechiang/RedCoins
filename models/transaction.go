package models

import (
	u "redcoins/utils"
	"github.com/jinzhu/gorm"
	"fmt"
)

type Transaction struct {
	gorm.Model
	Type string `json:"type"`
	Bitcoins float64 `json:"bitcoins"`
	Convert_Rt float64 `json:"convert_rt"`
	Final_Value float64 `json:"final_value"`
	User_Id_1 uint `json:"user_id_1"` //The user that this transaction belongs to
	User_Id_2 uint `json:"user_id_2"`
}

/*
 This struct function validate the required parameters sent through the http request body
returns message and true if the requirement is met
*/
func (transaction *Transaction) Validate() (map[string] interface{}, bool) {

	if transaction.Type == "" {
		return u.Message(false, "Transaction type should be on the payload"), false //implementar if type "sell" or "buy"
	}

	if transaction.Bitcoins <= 0 {
		return u.Message(false, "Bitcoins value should be on the payload"), false
	}

	if transaction.Convert_Rt <= 0 {
		return u.Message(false, "Convert Rate value should be on the payload"), false
	}

	if transaction.Final_Value <= 0 {
		return u.Message(false, "Final value should be on the payload"), false
	}

	if transaction.User_Id_1 <= 0 {
		return u.Message(false, "User 1 is not recognized"), false
	}

	if transaction.User_Id_2 <= 0 {
		return u.Message(false, "User 2 is not recognized"), false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (transaction *Transaction) Create() (map[string] interface{}) {

	if resp, ok := transaction.Validate(); !ok {
		return resp
	}

	GetDB().Create(transaction)

	resp := u.Message(true, "success")
	resp["transaction"] = transaction
	return resp
}

func GetTransaction(id uint) (*Transaction) {

	transaction := &Transaction{}
	err := GetDB().Table("transactions").Where("id = ?", id).First(transaction).Error
	if err != nil {
		return nil
	}
	return transaction
}

func GetTransactions(user uint) ([]*Transaction) {

	transactions := make([]*Transaction, 0)
	err := GetDB().Table("transactions").Where("user_id_1 = ?", user).Find(&transactions).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return transactions
}