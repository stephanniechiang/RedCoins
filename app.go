package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    // "os"
    "net/http"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

// const clientID = "3481f956589106d754c6"
// const clientSecret = "68179b61e46a6aa005a4cf8adfe809b3940a0af7"

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
    connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

    var err error
    a.DB, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()
    a.initializeRoutes()

    
}

func (a *App) Run(addr string) {

    log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
    // a.Router.Handle("/", http.FileServer(http.Dir("public")))
    
    // a.Router.HandleFunc("/oauth/redirect",a.auth)

    a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
    a.Router.HandleFunc("/user", a.createUser).Methods("POST")
    a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
    a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
    a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

    a.Router.HandleFunc("/transactions", a.getTransactions).Methods("GET")
    a.Router.HandleFunc("/transaction", a.createTransaction).Methods("POST")
    a.Router.HandleFunc("/transaction/{id:[0-9]+}", a.getTransaction).Methods("GET")
    a.Router.HandleFunc("/transaction/{id:[0-9]+}", a.updateTransaction).Methods("PUT")
    a.Router.HandleFunc("/transaction/{id:[0-9]+}", a.deleteTransaction).Methods("DELETE")
}

// func (a *App) auth(w http.ResponseWriter, r *http.Request) {
//         httpClient := http.Client{}
//         // First, we need to get the value of the `code` query param
//         err := r.ParseForm()
//         if err != nil {
//             fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
//             w.WriteHeader(http.StatusBadRequest)
//         }
//         code := r.FormValue("code")

//         // Next, lets for the HTTP request to call the github oauth enpoint
//         // to get our access token
//         reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
//         req, err := http.NewRequest(http.MethodPost, reqURL, nil)
//         if err != nil {
//             fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
//             w.WriteHeader(http.StatusBadRequest)
//         }
//         // We set this header since we want the response
//         // as JSON
//         req.Header.Set("accept", "application/json")

//         // Send out the HTTP request
//         res, err := httpClient.Do(req)
//         if err != nil {
//             fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
//             w.WriteHeader(http.StatusInternalServerError)
//         }
//         defer res.Body.Close()

//         // Parse the request body into the `OAuthAccessResponse` struct
//         var t OAuthAccessResponse
//         if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
//             fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
//             w.WriteHeader(http.StatusBadRequest)
//         }

//         // Finally, send a response to redirect the user to the "welcome" page
//         // with the access token
    
//         w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
//         w.WriteHeader(http.StatusFound)     
//     }

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    products, err := getUsers(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, products)
}

func (a *App) getTransactions(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    products, err := getTransactions(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
    var u user
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := u.createUser(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) createTransaction(w http.ResponseWriter, r *http.Request) {
    var u transaction
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := u.createTransaction(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    u := user{ID: id}
    if err := u.getUser(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "User not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getTransaction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    u := transaction{ID: id}
    if err := u.getTransaction(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Transaction not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    var u user
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    u.ID = id

    if err := u.updateUser(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateTransaction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    var u transaction
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    u.ID = id

    if err := u.updateTransaction(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid User ID")
        return
    }

    u := user{ID: id}
    if err := u.deleteUser(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteTransaction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Transaction ID")
        return
    }

    u := transaction{ID: id}
    if err := u.deleteTransaction(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}