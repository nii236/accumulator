package accumulator

import (
	"accumulator/db"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	vrc "github.com/nii236/vrchat-go/client"

	"net/http"
	"text/template"

	"github.com/caddyserver/caddy"
	// http driver for caddy
	_ "github.com/caddyserver/caddy/caddyhttp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"
)

// ErrNotImplemented is used to stub empty funcs
var ErrNotImplemented = errors.New("not implemented")

// HandlerFunc is a custom http.HandlerFunc that returns a status code and error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

// ErrorResponse for HTTP
type ErrorResponse struct {
	Err     string `json:"err"`
	Message string `json:"message"`
}

// Err constructor
func Err(err error, message ...string) *ErrorResponse {
	e := &ErrorResponse{
		Err: err.Error(),
	}
	e.Message = err.Error()
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

// Unwrap the inner error
func (e *ErrorResponse) Unwrap() error {
	return errors.New(e.Err)
}
func (e *ErrorResponse) Error() string {
	return e.Message
}

// JSON body for HTTP response
func (e *ErrorResponse) JSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
func withError(next HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		result, code, err := next(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, Err(err).JSON(), code)
			return
		}
		if result == nil {
			err := errors.New("no response")
			fmt.Println(err)
			http.Error(w, Err(err, "no response").JSON(), code)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			fmt.Println(err)
			http.Error(w, Err(err).JSON(), code)
			return
		}
		return
	}
	return fn
}

const caddyfileTemplate = `
{{ .caddyAddr}} {
	tls off
    proxy /api/ localhost{{ .apiAddr }} {
		transparent
		websocket
		timeout 10m
    }
    root {{ .rootPath }}
    rewrite { 
        if {path} not_match ^/api
        to {path} /
    }
}
`

// RunServer the service
func RunServer(ctx context.Context, conn *sqlx.DB, serverAddr string, log *zap.SugaredLogger) error {
	log.Infow("Start service", "svc-addr", serverAddr)
	// tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api", func(r chi.Router) {
		// Authenticated routes
		r.Group(func(r chi.Router) {
			// r.Use(jwtauth.Verifier(tokenAuth))
			// r.Use(jwtauth.Authenticator)

			r.Post("/auth/sign_out", withError(signOutHandler))
			r.Post("/auth/check", withError(checkHandler))
			r.Post("/auth/set_password", withError(setPasswordHandler))

			r.Get("/integrations/list", withError(integrationsListHandler))
			r.Post("/integrations/add_username", withError(integrationsAddUsernameHandler))
			r.Post("/integrations/delete", withError(integrationsDeleteHandler))

			r.Get("/friends/list", withError(friendListHandler))
			r.Post("/friends/refresh", withError(friendRefreshHandler))
			r.Post("/friends/promote", withError(friendPromoteHandler))

			r.Post("/teachers/demote", withError(teacherDemoteHandler))
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

	})

	return http.ListenAndServe(serverAddr, r)
}

// RunLoadBalancer starts Caddy
func RunLoadBalancer(ctx context.Context, conn *sqlx.DB, loadBalancerAddr, serverAddr, rootPath string, log *zap.SugaredLogger) error {
	log.Infow("Starting load balancer", "lb-addr", loadBalancerAddr, "svc-addr", serverAddr, "web", rootPath)
	caddy.AppName = "Accumulator"
	caddy.AppVersion = "0.0.1"
	caddy.Quiet = true
	t := template.Must(template.New("CaddyFile").Parse(caddyfileTemplate))
	data := map[string]string{
		"caddyAddr": loadBalancerAddr,
		"apiAddr":   serverAddr,
		"rootPath":  rootPath,
	}

	result := &bytes.Buffer{}
	err := t.Execute(result, data)
	if err != nil {
		return err
	}
	caddyfile := &caddy.CaddyfileInput{
		Contents:       result.Bytes(),
		Filepath:       "Caddyfile",
		ServerTypeName: "http",
	}

	instance, err := caddy.Start(caddyfile)
	if err != nil {
		return err
	}
	instance.Wait()
	return nil
}

func checkHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, 500, ErrNotImplemented
}
func integrationsListHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Data db.IntegrationSlice `json:"data"`
	}
	result, err := db.Integrations().AllG()
	if err != nil {
		return nil, 500, err
	}
	if result == nil {
		return &Response{}, 500, err
	}
	return &Response{result}, 200, nil
}
func integrationsAddUsernameHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type Response struct {
		Data string `json:"data"`
	}

	req := &Request{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	_, apiKey, authToken, err := vrc.Token(vrc.ReleaseAPIURL, req.Username, req.Password)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	record := &db.Integration{
		APIKey:    apiKey,
		AuthToken: authToken,
	}
	client, err := vrc.NewClient(vrc.ReleaseAPIURL, authToken, apiKey)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	client.RemoteConfig()

	// TODO: Figure out how to handle error handling with ROWID primary keys
	record.InsertG(boil.Infer())
	return record, http.StatusOK, nil
}
func integrationsDeleteHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, 500, ErrNotImplemented
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

func friendListHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Request struct {
		IntegrationID int
	}
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	req := &Request{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = RefreshFriendCache(req.IntegrationID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
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
func teacherDemoteHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
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
