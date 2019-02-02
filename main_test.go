package main

import (
    "fmt"
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "strconv"
    "testing"
)

var a App

func TestMain(m *testing.M) {
    a = App{}
    a.Initialize("root", "moedasvermelhas", "RedCoins")

    ensureTableExistsUser()
    ensureTableExistsTransaction()

    code := m.Run()

    clearTableUser()
    clearTableTransaction()

    os.Exit(code)
}

func ensureTableExistsUser() {
    if _, err := a.DB.Exec(tableCreationQueryUser); err != nil {
        log.Fatal(err)
    }
}

func ensureTableExistsTransaction() {
    if _, err := a.DB.Exec(tableCreationQueryTransaction); err != nil {
        log.Fatal(err)
    }
}

func clearTableUser() {
    a.DB.Exec("DELETE FROM users")
    a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func clearTableTransaction() {
    a.DB.Exec("DELETE FROM transactions")
    a.DB.Exec("ALTER TABLE transactions AUTO_INCREMENT = 1")
}

const tableCreationQueryUser = `
CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    birthday VARCHAR(50) NOT NULL,
    balance FLOAT,
    balance_bit FLOAT
);`

const tableCreationQueryTransaction = `
CREATE TABLE IF NOT EXISTS transactions
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR(50) NOT NULL,
    hour VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    bitcoins FLOAT NOT NULL,
    convert_tx FLOAT NOT NULL,
    final_value FLOAT NOT NULL,
    user_id_1 INT NOT NULL,
    user_id_2 INT NOT NULL,
    FOREIGN KEY (user_id_1) REFERENCES users(id)
);`

func TestEmptyTableUser(t *testing.T) {
    clearTableUser()

    req, _ := http.NewRequest("GET", "/users", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func TestEmptyTableTransaction(t *testing.T) {
    clearTableTransaction()

    req, _ := http.NewRequest("GET", "/transactions", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestGetNonExistentUser(t *testing.T) {
    clearTableUser()

    req, _ := http.NewRequest("GET", "/user/45", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "User not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
    }
}

func TestGetNonExistentTransaction(t *testing.T) {
    clearTableTransaction()

    req, _ := http.NewRequest("GET", "/transaction/45", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Transaction not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Transaction not found'. Got '%s'", m["error"])
    }
}

func TestCreateUser(t *testing.T) {
    clearTableUser()

    payload := []byte(`{"name":"John","last_name":"Appleseed","email":"john@appleseed.com","password":"mortadela1","birthday":"01/01/1900","balance":0.00,"balance_bit":0.00}`)

    req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
    response := executeRequest(req)

    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["name"] != "John" {
        t.Errorf("Expected user name to be 'John'. Got '%v'", m["name"])
    }

    if m["last_name"] != "Appleseed" {
        t.Errorf("Expected user last_name to be 'Appleseed'. Got '%v'", m["last_name"])
    }

    if m["email"] != "john@appleseed.com" {
        t.Errorf("Expected user email to be 'john@appleseed.com'. Got '%v'", m["email"])
    }

    if m["password"] != "mortadela1" {
        t.Errorf("Expected user password to be 'mortadela1'. Got '%v'", m["password"])
    }

    if m["birthday"] != "01/01/1900" {
        t.Errorf("Expected user birthday to be '01/01/1900'. Got '%v'", m["birthday"])
    }

    if m["balance"] != 0.00 {
        t.Errorf("Expected user balance to be '0.00'. Got '%v'", m["balance"])
    }

    if m["balance_bit"] != 0.00 {
        t.Errorf("Expected user balance_bit to be '0.00'. Got '%v'", m["balance_bit"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
    }
}

func TestCreateTransaction(t *testing.T) {
    clearTableTransaction()

    // clearTableUser()

    // payload1 := []byte(`{"name":"John","last_name":"Appleseed","email":"john@appleseed.com","password":"mortadela1","birthday":"01/01/1900","balance":0.00,"balance_bit":0.00}`)

    // req1, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload1))
    // response1 := executeRequest(req1)

    // checkResponseCode(t, http.StatusCreated, response1.Code)
    addUsers(1);

    payload2 := []byte(`{"date":"03/04/2015","hour":"18:40","type":"sale","bitcoins":100.00,"convert_tx":2.30,"final_value":230.00,"user_id_1":1,"user_id_2":43}`)

    req2, _ := http.NewRequest("POST", "/transaction", bytes.NewBuffer(payload2))
    response2 := executeRequest(req2)

    checkResponseCode(t, http.StatusCreated, response2.Code)

    var m map[string]interface{}
    json.Unmarshal(response2.Body.Bytes(), &m)

    if m["date"] != "03/04/2015" {
        t.Errorf("Expected transaction date to be '03/04/2015'. Got '%v'", m["date"])
    }

    if m["hour"] != "18:40" {
        t.Errorf("Expected transaction hour to be '18:40'. Got '%v'", m["hour"])
    }

    if m["type"] != "sale" {
        t.Errorf("Expected transaction type to be 'sale'. Got '%v'", m["type"])
    }

    if m["bitcoins"] != 100.00 {
        t.Errorf("Expected transaction bitcoins to be '100.00'. Got '%v'", m["bitcoins"])
    }

    if m["convert_tx"] != 2.30 {
        t.Errorf("Expected transaction convert_tx to be '2.30'. Got '%v'", m["convert_tx"])
    }

    if m["final_value"] != 230.00 {
        t.Errorf("Expected transaction final_value to be '230.00'. Got '%v'", m["final_value"])
    }

    if m["user_id_1"] != 1 {
        t.Errorf("Expected transaction user_id_1 to be '1'. Got '%v'", m["user_id_1"])
    }

    if m["user_id_2"] != 43 {
        t.Errorf("Expected transaction user_id_2 to be '43'. Got '%v'", m["user_id_2"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected transaction ID to be '1'. Got '%v'", m["id"])
    }
}

func addUsers(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
        statement := fmt.Sprintf("INSERT INTO users(name, last_name, email, password, birthday, balance, balance_bit) VALUES('%s', '%s', '%s', '%s', '%s', %d, %d)", ("Name " + strconv.Itoa(i+1)), ("Last Name " + strconv.Itoa(i+1)), ("Email " + strconv.Itoa(i+1)), ("Password " + strconv.Itoa(i+1)), ("Birthday " + strconv.Itoa(i+1)), ((i+1.0) * 10.0), ((i+1.0) * 10.0))
        a.DB.Exec(statement)
    }
}

func addTransactions(count int) {
    addUsers(1);
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
        statement := fmt.Sprintf("INSERT INTO transactions(date, hour, type, bitcoins, convert_tx, final_value, user_id_1, user_id_2) VALUES('%s', '%s', '%s', %d, %d, %d, %d, %d)", ("Date " + strconv.Itoa(i+1)), ("Hour " + strconv.Itoa(i+1)), ("Type " + strconv.Itoa(i+1)), ((i+1.0) * 10.0), ((i+1.0) * 10.0), ((i+1.0) * 10.0), 1, ((i+1.0) * 10.0))
        a.DB.Exec(statement)
    }
}

func TestGetUser(t *testing.T) {
    clearTableUser()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTransaction(t *testing.T) {
    clearTableTransaction()
    addTransactions(1)

    req, _ := http.NewRequest("GET", "/transaction/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
    clearTableUser()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)
    var originalUser map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalUser)

    payload := []byte(`{"name":"John - Updated","last_name":"Appleseed","email":"john@appleseed.com","password":"mortadela1","birthday":"01/01/1900","balance":0.00,"balance_bit":0.00}`)

    req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["id"] != originalUser["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
    }

    if m["name"] == originalUser["name"] {
        t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalUser["name"], m["name"], m["name"])
    }

    if m["last_name"] == originalUser["last_name"] {
        t.Errorf("Expected the last_name to change from '%v' to '%v'. Got '%v'", originalUser["last_name"], m["last_name"], m["last_name"])
    }

    if m["email"] == originalUser["email"] {
        t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", originalUser["email"], m["email"], m["email"])
    }

    if m["password"] == originalUser["password"] {
        t.Errorf("Expected the password to change from '%v' to '%v'. Got '%v'", originalUser["password"], m["password"], m["password"])
    }

    if m["birthday"] == originalUser["birthday"] {
        t.Errorf("Expected the birthday to change from '%v' to '%v'. Got '%v'", originalUser["birthday"], m["birthday"], m["birthday"])
    }

    if m["balance"] == originalUser["balance"] {
        t.Errorf("Expected the balance to change from '%v' to '%v'. Got '%v'", originalUser["balance"], m["balance"], m["balance"])
    }

    if m["balance_bit"] == originalUser["balance_bit"] {
        t.Errorf("Expected the balance_bit to change from '%v' to '%v'. Got '%v'", originalUser["balance_bit"], m["balance_bit"], m["balance_bit"])
    }
}

func TestUpdateTransaction(t *testing.T) {
    clearTableTransaction()
    addTransactions(1)

    req, _ := http.NewRequest("GET", "/transaction/1", nil)
    response := executeRequest(req)
    var originalTransaction map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalTransaction)

    payload := []byte(`{"date":"31/01/2019","hour":"16:50","type":"sale","bitcoins":10.00,"convert_tx":2.3,"final_value":23.00,"user_id_1":2,"user_id_2":43}`)

    req, _ = http.NewRequest("PUT", "/transaction/1", bytes.NewBuffer(payload))
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["id"] != originalTransaction["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalTransaction["id"], m["id"])
    }

    if m["date"] == originalTransaction["date"] {
        t.Errorf("Expected the date to change from '%v' to '%v'. Got '%v'", originalTransaction["date"], m["date"], m["date"])
    }

    if m["hour"] == originalTransaction["hour"] {
        t.Errorf("Expected the hour to change from '%v' to '%v'. Got '%v'", originalTransaction["hour"], m["hour"], m["hour"])
    }

    if m["type"] == originalTransaction["type"] {
        t.Errorf("Expected the type to change from '%v' to '%v'. Got '%v'", originalTransaction["type"], m["type"], m["type"])
    }

    if m["bitcoins"] == originalTransaction["bitcoins"] {
        t.Errorf("Expected the bitcoins to change from '%v' to '%v'. Got '%v'", originalTransaction["bitcoins"], m["bitcoins"], m["bitcoins"])
    }

    if m["convert_tx"] == originalTransaction["convert_tx"] {
        t.Errorf("Expected the convert_tx to change from '%v' to '%v'. Got '%v'", originalTransaction["convert_tx"], m["convert_tx"], m["convert_tx"])
    }

    if m["final_value"] == originalTransaction["final_value"] {
        t.Errorf("Expected the final_value to change from '%v' to '%v'. Got '%v'", originalTransaction["final_value"], m["final_value"], m["final_value"])
    }

    if m["user_id_1"] == originalTransaction["user_id_1"] {
        t.Errorf("Expected the user_id_1 to change from '%v' to '%v'. Got '%v'", originalTransaction["user_id_1"], m["user_id_1"], m["user_id_1"])
    }

    if m["user_id_2"] == originalTransaction["user_id_2"] {
        t.Errorf("Expected the user_id_2 to change from '%v' to '%v'. Got '%v'", originalTransaction["user_id_2"], m["user_id_2"], m["user_id_2"])
    }
}

func TestDeleteUser(t *testing.T) {
    clearTableUser()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/user/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/user/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeleteTransaction(t *testing.T) {
    clearTableTransaction()
    clearTableUser()
    addUsers(1)
    addTransactions(1)

    req, _ := http.NewRequest("GET", "/transaction/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/transaction/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/transaction/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}