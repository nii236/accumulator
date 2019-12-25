package main

import (
	"accumulator/db"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/boil"
	"log"
	"net/http"
)

// ErrNotImplemented is used to stub empty funcs
var ErrNotImplemented = errors.New("not implemented")

// HandlerFunc is a custom http.HandlerFunc that returns a status code and error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

func withError(next HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		result, code, err := next(w, r)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		if result == nil {
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		return
	}
	return fn
}

func main() {
	addr := flag.String("addr", ":8080", "Address to host on")
	flag.Parse()
	conn, err := connect()
	if err != nil {
		fmt.Println(err)
	}
	boil.SetDB(conn)
	// tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		// r.Use(jwtauth.Verifier(tokenAuth))
		// r.Use(jwtauth.Authenticator)

		r.Post("/auth/sign_out", withError(signOutHandler))
		r.Post("/auth/set_password", withError(setPasswordHandler))
		r.Post("/auth/set_token", withError(setTokenHandler))

		r.Get("/friends/list", withError(friendListHandler))
		r.Post("/friends/refresh", withError(friendRefreshHandler))
		r.Post("/friends/promote", withError(friendPromoteHandler))
		r.Post("/friends/demote", withError(friendDemoteHandler))

		r.Get("/teachers/list", withError(teacherListHandler))
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/sign_in", withError(signInHandler))
		r.Post("/auth/sign_up", withError(signUpHandler))
		r.Get("/auth/forgot_password", withError(signUpHandler))
		r.Post("/auth/request_password_reset", withError(signUpHandler))
		r.Post("/auth/reset_password", withError(signUpHandler))
		r.Get("/metrics", withError(metricsHandler))
	})

	fmt.Println("Running accumulator on", *addr)
	log.Fatalln(http.ListenAndServe(*addr, r))
}

func connect() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("sqlite3", "./accumulator.db")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func storeToken(username, authToken, apiKey string) error {
	return nil
}
func loadToken(username string) (string, string, error) {
	return "", "", nil
}

func signOutHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, 200, ErrNotImplemented
}
func signInHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Request struct {
		Email    string
		Password string
	}
	type Response struct {
		Data string `json:"data"`
	}
	return nil, 200, ErrNotImplemented
}
func signUpHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Request struct {
		Email    string
		Password string
	}
	type Response struct {
		Data string `json:"data"`
	}
	return nil, 200, ErrNotImplemented
}
func setPasswordHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Request struct {
		Password string
	}
	return nil, 200, ErrNotImplemented
}
func setTokenHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Request struct {
		APIKey    string
		AuthToken string
	}
	return nil, 200, ErrNotImplemented
}
func friendListHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	result, err := db.Friends(db.FriendWhere.IsTeacher.EQ(false)).AllG()
	if err != nil {
		return nil, 500, err
	}
	if result == nil {
		return &Response{}, 500, err
	}
	return &Response{result}, 200, nil
}
func friendRefreshHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	// TODO: Manually refresh friend locations
	result, err := db.Friends().AllG()
	if err != nil {
		return nil, 500, err
	}
	if result == nil {
		return &Response{}, 500, err
	}
	return &Response{result}, 200, nil
}
func friendPromoteHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Data *db.Friend `json:"data"`
	}
	return nil, 200, ErrNotImplemented
}
func friendDemoteHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Data *db.Friend `json:"data"`
	}
	return nil, 200, ErrNotImplemented
}
func teacherListHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	result, err := db.Friends(db.FriendWhere.IsTeacher.EQ(true)).AllG()
	if err != nil {
		return nil, 500, err
	}
	if result == nil {
		return &Response{}, 500, err
	}
	return &Response{result}, 200, nil
}

func metricsHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, 200, ErrNotImplemented
}
