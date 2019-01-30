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

    ensureTableExists()

    code := m.Run()

    clearTable()

    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("DELETE FROM users")
    a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    birthday VARCHAR(50) NOT NULL,
    balance FLOAT NOT NULL
)`

func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/users", nil)
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
    clearTable()

    req, _ := http.NewRequest("GET", "/user/45", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "User not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
    }
}

func TestCreateUser(t *testing.T) {
    clearTable()

    payload := []byte(`{"name":"John","last_name":"Appleseed","email":"john@appleseed.com","password":"mortadela1","birthday":"01/01/1900","balance":0.00}`)

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

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
    }
}

func addUsers(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
    statement := fmt.Sprintf("INSERT INTO users(name, last_name, email, password, birthday, balance) VALUES('%s', '%s', '%s', '%s', '%s', %d)", ("User " + strconv.Itoa(i+1)), ("User " + strconv.Itoa(i+1)), ("User " + strconv.Itoa(i+1)), ("User " + strconv.Itoa(i+1)), ("User " + strconv.Itoa(i+1)), ((i+1.0) * 10.0))
        a.DB.Exec(statement)
    }
}

func TestGetUser(t *testing.T) {
    clearTable()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
    clearTable()
    addUsers(1)

    req, _ := http.NewRequest("GET", "/user/1", nil)
    response := executeRequest(req)
    var originalUser map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalUser)

    payload := []byte(`{"name":"John - Updated","last_name":"Appleseed","email":"john@appleseed.com","password":"mortadela1","birthday":"01/01/1900","balance":0.00}`)

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
}

func TestDeleteUser(t *testing.T) {
    clearTable()
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