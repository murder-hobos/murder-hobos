package routes

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func newLoginRouter(env *Env) *mux.Router {
	stdChain := alice.New(env.withClaims)
	r := mux.NewRouter()

	r.Handle("/login", stdChain.ThenFunc(env.loginIndex)).Methods("GET")
	r.Handle("/login", stdChain.ThenFunc(env.loginProcess)).Methods("POST")
	// r.HandleFunc("/login/register", env.loginRegister).Methods("POST")
	r.Handle("/logout", stdChain.ThenFunc(env.logoutProcess))

	return r
}

// Claims is a struct that represents the data stored in our
// tokens
type Claims struct {
	UID      int    `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// errors will be rather bare since we're only doing
// username/password
type loginPage struct {
	Errors []string
}

func (env *Env) loginIndex(w http.ResponseWriter, r *http.Request) {
	if u := r.Context().Value("User"); u != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	lp := r.Context().Value("LoginPage")

	if t, ok := tmpls["login.html"]; ok {
		t.ExecuteTemplate(w, "base", lp)
	} else {
		errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (env *Env) loginProcess(w http.ResponseWriter, r *http.Request) {
	u := r.PostFormValue("username")
	p := r.PostFormValue("password")

	if u == "" || p == "" {
		loginPageWithErrors(w, r, "Username or Password cannot be empty.")
		return
	}

	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	//Authenticate user
	user, ok := env.db.GetUserByUsername(u)
	if !ok {
		loginPageWithErrors(w, r, "Invalid Username or Password")
		return
	}

	// TODO: CHANGE THIS TO BCRYPT AFTER DONE TESTING
	if bytes.Equal(user.Password, []byte(p)) {
		claims := Claims{
			user.ID,
			user.Username,
			jwt.StandardClaims{
				ExpiresAt: expireToken,
				Issuer:    "localhost:8081",
			},
		}

		// Generate signed token with our claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SIGNING_KEY")))
		if err != nil {
			http.Error(w, "Error creating signed token", http.StatusInternalServerError)
		}

		// Store our token on client
		cookie := http.Cookie{
			Name:     "Auth",
			Value:    signedToken,
			Expires:  expireCookie,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Invalid password
	loginPageWithErrors(w, r, "Invalid Username or Password")
	return
}

func (env *Env) logoutProcess(w http.ResponseWriter, r *http.Request) {
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(w, &deleteCookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// little utility
func loginPageWithErrors(w http.ResponseWriter, r *http.Request, errs ...string) {
	lp := loginPage{Errors: errs}
	ctx := context.WithValue(r.Context(), "LoginPage", &lp)
	r = r.WithContext(ctx)
	r.Method = "GET"
	http.Redirect(w, r, "/login", http.StatusBadRequest)
	return

}
