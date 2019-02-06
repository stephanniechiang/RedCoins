package main

import (
    "net/http"
    // "net/url"
    // "log"
    "fmt"
    "os"
    "encoding/json"

 //    "github.com/go-session/session"
	// "gopkg.in/oauth2.v3/errors"
	// "gopkg.in/oauth2.v3/manage"
	// "gopkg.in/oauth2.v3/models"
	// "gopkg.in/oauth2.v3/server"
	// "gopkg.in/oauth2.v3/store"
)

// const clientID = "3481f956589106d754c6"
// const clientSecret = "68179b61e46a6aa005a4cf8adfe809b3940a0af7"

func main() {
	a := App{} 
	a.Initialize("root", "moedasvermelhas", "RedCoins")
	a.Run(":8080")
	http.ListenAndServe(":8080", nil)

	// fs := http.FileServer(http.Dir("public"))
	// http.Handle("/", fs)

	// // We will be using `httpClient` to make external HTTP requests later in our code
	// httpClient := http.Client{}

	// // Create a new redirect route route
	// http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
	// 	// First, we need to get the value of the `code` query param
	// 	err := r.ParseForm()
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	}
	// 	code := r.FormValue("code")

	// 	// Next, lets for the HTTP request to call the github oauth enpoint
	// 	// to get our access token
	// 	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	// 	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	}
	// 	// We set this header since we want the response
	// 	// as JSON
	// 	req.Header.Set("accept", "application/json")

	// 	// Send out the HTTP request
	// 	res, err := httpClient.Do(req)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 	}
	// 	defer res.Body.Close()

	// 	// Parse the request body into the `OAuthAccessResponse` struct
	// 	var t OAuthAccessResponse
	// 	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
	// 		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	}

	// 	// Finally, send a response to redirect the user to the "welcome" page
	// 	// with the access token
	// 	w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
	// 	w.WriteHeader(http.StatusFound)
	// })

	// http.ListenAndServe(":8080", nil)
}

// type OAuthAccessResponse struct {
// 	AccessToken string `json:"access_token"`
// }





// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"time"

// 	"github.com/go-session/session"
// 	"gopkg.in/oauth2.v3/errors"
// 	"gopkg.in/oauth2.v3/manage"
// 	"gopkg.in/oauth2.v3/models"
// 	"gopkg.in/oauth2.v3/server"
// 	"gopkg.in/oauth2.v3/store"
// )

// func main() {
// 	manager := manage.NewDefaultManager()
// 	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

// 	// token store
// 	manager.MustTokenStorage(store.NewMemoryTokenStore())

// 	clientStore := store.NewClientStore()
// 	clientStore.Set("222222", &models.Client{
// 		ID:     "222222",
// 		Secret: "22222222",
// 		Domain: "http://localhost:9094",
// 	})
// 	manager.MapClientStorage(clientStore)

// 	srv := server.NewServer(server.NewConfig(), manager)
// 	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

// 	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
// 		log.Println("Internal Error:", err.Error())
// 		return
// 	})

// 	srv.SetResponseErrorHandler(func(re *errors.Response) {
// 		log.Println("Response Error:", re.Error.Error())
// 	})

// 	http.HandleFunc("/login", loginHandler)
// 	http.HandleFunc("/auth", authHandler)

// 	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
// 		store, err := session.Start(nil, w, r)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		var form url.Values
// 		if v, ok := store.Get("ReturnUri"); ok {
// 			form = v.(url.Values)
// 		}
// 		r.Form = form

// 		store.Delete("ReturnUri")
// 		store.Save()

// 		err = srv.HandleAuthorizeRequest(w, r)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 		}
// 	})

// 	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
// 		err := srv.HandleTokenRequest(w, r)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	})

// 	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
// 		token, err := srv.ValidationBearerToken(r)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		data := map[string]interface{}{
// 			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
// 			"client_id":  token.GetClientID(),
// 			"user_id":    token.GetUserID(),
// 		}
// 		e := json.NewEncoder(w)
// 		e.SetIndent("", "  ")
// 		e.Encode(data)
// 	})

// 	log.Println("Server is running at 9096 port.")
// 	log.Fatal(http.ListenAndServe(":9096", nil))
// }

// func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
// 	store, err := session.Start(nil, w, r)
// 	if err != nil {
// 		return
// 	}

// 	uid, ok := store.Get("LoggedInUserID")
// 	if !ok {
// 		if r.Form == nil {
// 			r.ParseForm()
// 		}

// 		store.Set("ReturnUri", r.Form)
// 		store.Save()

// 		w.Header().Set("Location", "/login")
// 		w.WriteHeader(http.StatusFound)
// 		return
// 	}

// 	userID = uid.(string)
// 	store.Delete("LoggedInUserID")
// 	store.Save()
// 	return
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	store, err := session.Start(nil, w, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if r.Method == "POST" {
// 		store.Set("LoggedInUserID", "000000")
// 		store.Save()

// 		w.Header().Set("Location", "/auth")
// 		w.WriteHeader(http.StatusFound)
// 		return
// 	}
// 	outputHTML(w, r, "static/login.html")
// }

// func authHandler(w http.ResponseWriter, r *http.Request) {
// 	store, err := session.Start(nil, w, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if _, ok := store.Get("LoggedInUserID"); !ok {
// 		w.Header().Set("Location", "/login")
// 		w.WriteHeader(http.StatusFound)
// 		return
// 	}

// 	outputHTML(w, r, "static/auth.html")
// }

// func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	defer file.Close()
// 	fi, _ := file.Stat()
// 	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
// }