package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Claims is a struct that represents the data stored in our
// tokens
type Claims struct {
	UID      int    `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (env *Env) loginIndex(w http.ResponseWriter, r *http.Request) {
	// already authenticated
	if claims := r.Context().Value("Claims"); claims != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	errs := r.Context().Value("Errors")
	log.Println(errs)
	// We keep claims here because base template requires at least a nil claims
	data := map[string]interface{}{
		"Claims": nil,
		"Errors": errs,
	}

	if t, ok := tmpls["login.html"]; ok {
		t.ExecuteTemplate(w, "base", data)
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

	//Authenticate user
	user, ok := env.db.GetUserByUsername(u)
	if !ok {
		loginPageWithErrors(w, r, "Invalid Username or Password")
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(p)); err != nil {
		// Invalid password
		loginPageWithErrors(w, r, "Invalid Username or Password")
		return

	}
	assignToken(w, user.ID, user.Username)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

// Logs a user out by "deleting" their token
func (env *Env) logoutProcess(w http.ResponseWriter, r *http.Request) {
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(w, &deleteCookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Creates a new user
func (env *Env) registerProcess(w http.ResponseWriter, r *http.Request) {
	u := r.PostFormValue("username")
	p := r.PostFormValue("password")
	pc := r.PostFormValue("confirm-password")

	if u == "" || p == "" || pc == "" {
		loginPageWithErrors(w, r, "Username or Password cannot be empty.")
		return
	}

	if p != pc {
		loginPageWithErrors(w, r, "Passwords do not match.")
	}

	if user, ok := env.db.CreateUser(u, p); ok {
		log.Println("User created")
		assignToken(w, user.ID, user.Username)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	errorHandler(w, r, http.StatusInternalServerError)
}

// assignToken stores a cookie with our auth token
func assignToken(w http.ResponseWriter, id int, uname string) {
	expireToken := time.Now().Add(time.Hour * 8).Unix()
	expireCookie := time.Now().Add(time.Hour * 8)

	claims := Claims{
		id,
		uname,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "murder-hobos",
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

}

// little utility
func loginPageWithErrors(w http.ResponseWriter, r *http.Request, errs ...string) {
	r.Method = "GET"
	ctx := context.WithValue(r.Context(), "Errors", errs)
	http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
	return
}
