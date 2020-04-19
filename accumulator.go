package accumulator

import (
	"accumulator/db"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dgrijalva/jwt-go"
	vrc "github.com/nii236/vrchat-go/client"

	"net/http"
	"text/template"

	"github.com/caddyserver/caddy"
	// http driver for caddy
	"github.com/alexedwards/scs/v2"
	_ "github.com/caddyserver/caddy/caddyhttp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"
)

var sessionManager *scs.SessionManager

// ErrNotImplemented is used to stub empty funcs
var ErrNotImplemented = errors.New("not implemented")

// ErrUnableToPopulate occurs because of SQLite's ID creation order
var ErrUnableToPopulate = "db: unable to populate default values"

// SecureHandlerFunc is a custom http.HandlerFunc that returns a status code and error
type SecureHandlerFunc func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error)

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
func withUser(auther *Auther, next SecureHandlerFunc) HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		jwtString := ""
		cookie, err := r.Cookie("jwt")
		if err != nil {
			jwtString = strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
		}
		if cookie != nil {
			jwtString = cookie.Value
		}
		if jwtString == "" {
			return nil, http.StatusUnauthorized, errors.New("no jwt provided in cookie or header")
		}

		token, err := auther.TokenAuth.Decode(jwtString)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("could not cast to jwt.MapClaims")
		}
		idI, ok := claims["id"]
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("could not read value for key id")
		}
		idStr, ok := idI.(string)
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("could not cast id to string")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		u, err := db.FindUserG(null.Int64From(int64(id)))
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return next(w, r, u)
	}
	return fn
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
			if err == nil {
				err = errors.New("no response")
			}
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
func RunServer(ctx context.Context, conn *sqlx.DB, serverAddr string, jwtsecret string, d *Darer, log *zap.SugaredLogger) error {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	log.Infow("start api", "svc-addr", serverAddr)
	auther := NewAuther(jwtsecret)
	c := &API{log}

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api", func(r chi.Router) {
		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/sign_out", withError(c.signOutHandler))
			r.Get("/auth/check", withError(withUser(auther, c.checkHandler(auther))))
			r.Post("/auth/set_password", withError(withUser(auther, c.setPasswordHandler())))
			r.Get("/auth/jwt", withError(withUser(auther, c.userJWTHandler(auther))))

			r.Get("/blobs/{blob_id}", c.blobHandler())

			r.Get("/users/list", withError(withUser(auther, c.userListHandler())))
			r.Post("/users/impersonate/{user_id}", withError(withUser(auther, c.userImpersonateHandler(auther))))

			r.Get("/integrations/list", withError(withUser(auther, c.integrationsListHandler)))
			r.Post("/integrations/add_username", withError(withUser(auther, c.integrationsAddUsernameHandler(d))))
			r.Post("/integrations/{integration_id}/update_friends", withError(withUser(auther, c.integrationUpdateFriendsHandler(d))))
			r.Post("/integrations/{integration_id}/delete", withError(withUser(auther, c.integrationsDeleteHandler)))
			r.Get("/integrations/{integration_id}/attendance/{teacher_id}/list", withError(withUser(auther, c.attendanceListHandler)))
			r.Get("/integrations/{integration_id}/friends/list", withError(withUser(auther, c.friendListHandler)))
			r.Post("/integrations/{integration_id}/friends/refresh", withError(withUser(auther, c.friendRefreshHandler)))
			r.Post("/integrations/{integration_id}/friends/{friend_id}/promote", withError(withUser(auther, c.friendPromoteHandler)))
			r.Post("/integrations/{integration_id}/friends/{friend_id}/demote", withError(withUser(auther, c.friendDemoteHandler)))
		})

		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/sign_in", withError(c.signInHandler(auther)))
			r.Post("/auth/sign_up", withError(c.signUpHandler(auther)))
			r.Get("/metrics", promhttp.Handler().ServeHTTP)
		})

	})

	return http.ListenAndServe(serverAddr, sessionManager.LoadAndSave(r))
}

type API struct {
	log *zap.SugaredLogger
}

// RunLoadBalancer starts Caddy
func RunLoadBalancer(ctx context.Context, conn *sqlx.DB, loadBalancerAddr, serverAddr, rootPath string, log *zap.SugaredLogger) error {
	log.Infow("start load balancer", "lb-addr", loadBalancerAddr, "svc-addr", serverAddr, "web", rootPath)
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

func isIntegrationOwner(IntegrationID int, userID null.Int64) error {
	integration, err := db.FindIntegrationG(null.Int64From(int64(IntegrationID)))
	if err != nil {
		return err
	}
	if userID.Int64 != integration.UserID {
		return errors.New("unauthorized")
	}
	return nil
}

func (c *API) integrationUpdateFriendsHandler(d *Darer) func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		IntegrationIDStr := chi.URLParam(r, "integration_id")
		type Response struct {
			Success bool `json:"success,omitempty"`
		}
		IntegrationID, err := strconv.Atoi(IntegrationIDStr)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = isIntegrationOwner(IntegrationID, u.ID)
		if err != nil {
			return nil, http.StatusForbidden, err
		}
		err = refreshFriendCache(d, IntegrationID, true)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return &Response{true}, http.StatusOK, nil
	}
	return fn
}
func (c *API) checkHandler(auther *Auther) func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		type Response struct {
			Data *db.User `json:"data"`
		}
		cookie, err := r.Cookie("jwt")
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		_, err = auther.TokenAuth.Decode(cookie.Value)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		u.PasswordHash = ""
		return &Response{u}, 200, nil
	}
	return fn

}

func (c *API) integrationsListHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	type Response struct {
		Data db.IntegrationSlice `json:"data"`
	}
	result, err := db.Integrations(db.IntegrationWhere.UserID.EQ(u.ID.Int64)).AllG()
	if err != nil {
		return nil, 500, err
	}
	if result == nil {
		return &Response{}, 500, err
	}
	return &Response{result}, 200, nil
}
func (c *API) integrationsAddUsernameHandler(d *Darer) func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
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

		encryptedAuthToken, nonce, err := d.encrypt([]byte(authToken))
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		record := &db.Integration{
			UserID:         u.ID.Int64,
			Username:       req.Username,
			APIKey:         apiKey,
			AuthToken:      encryptedAuthToken,
			AuthTokenNonce: nonce,
		}

		exists, err := db.Integrations(db.IntegrationWhere.Username.EQ(req.Username)).ExistsG()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if exists {
			c.log.Infow("integration email exists, updating...", "email", req.Username)
			existingRecord, err := db.Integrations(db.IntegrationWhere.Username.EQ(req.Username)).OneG()
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			existingRecord.UserID = u.ID.Int64
			existingRecord.APIKey = apiKey

			encryptedAuthToken, nonce, err := d.encrypt([]byte(authToken))
			if err != nil {
				return nil, http.StatusBadRequest, err
			}
			existingRecord.AuthToken = encryptedAuthToken
			existingRecord.AuthTokenNonce = nonce
			_, err = existingRecord.UpdateG(boil.Whitelist(
				db.IntegrationColumns.UserID,
				db.IntegrationColumns.Username,
				db.IntegrationColumns.APIKey,
				db.IntegrationColumns.AuthToken,
			))
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			return record, http.StatusOK, nil
		}
		c.log.Infow("integration email does not exist, creating...", "email", req.Username)
		err = record.InsertG(boil.Infer())
		if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
			return nil, http.StatusInternalServerError, err
		}
		return record, http.StatusOK, nil
	}
	return fn
}
func (c *API) integrationsDeleteHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	IntegrationIDStr := chi.URLParam(r, "integration_id")
	IntegrationID, err := strconv.Atoi(IntegrationIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = isIntegrationOwner(IntegrationID, u.ID)
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	integration, err := db.FindIntegrationG(null.Int64From(int64(IntegrationID)))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	_, err = integration.DeleteG()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return nil, 200, nil
}

func (c *API) signOutHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	type Response struct {
		Success bool `json:"success"`
	}
	cookie := http.Cookie{Name: "jwt", Value: "", Expires: time.Unix(0, 0), HttpOnly: true, Path: "/", SameSite: http.SameSiteDefaultMode, Secure: false}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	return &Response{true}, http.StatusOK, nil
}
func (c *API) signInHandler(auther *Auther) func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		type Request struct {
			Email    string
			Password string
		}
		type Response struct {
			Success bool `json:"success"`
		}
		failedMessage := errors.New("Bad username or password")
		req := &Request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		defer r.Body.Close()

		err = auther.ValidatePassword(req.Email, req.Password)
		if err != nil {
			return nil, http.StatusBadRequest, failedMessage
		}

		user, err := db.Users(db.UserWhere.Email.EQ(req.Email)).OneG()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		expiration := time.Now().Add(time.Duration(30) * time.Hour * 24)
		jwt, err := auther.GenerateJWT(user.Email, strconv.Itoa(int(user.ID.Int64)), user.Role, expiration)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		cookie := http.Cookie{Name: "jwt", Value: jwt, Expires: expiration, HttpOnly: true, Path: "/", SameSite: http.SameSiteDefaultMode, Secure: false}
		http.SetCookie(w, &cookie)

		return &Response{true}, http.StatusOK, nil
	}
	return fn
}
func (c *API) signUpHandler(auther *Auther) func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		type Request struct {
			Email    string
			Password string
		}
		type Response struct {
			Data string `json:"data"`
		}
		req := &Request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		u := &db.User{
			Email:        req.Email,
			PasswordHash: HashPassword(req.Password),
		}
		err = u.InsertG(boil.Infer())
		if err != nil && !strings.Contains(err.Error(), ErrUnableToPopulate) {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 200, ErrNotImplemented
	}
	return fn
}
func (c *API) blobHandler() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		blobFilename := chi.URLParam(r, "blob_id")
		blob, err := db.Blobs(db.BlobWhere.FileName.EQ(blobFilename)).OneG()
		if err != nil {
			http.Error(w, Err(err).JSON(), http.StatusBadRequest)
			return
		}

		// tell the browser the returned content should be downloaded/inline
		if blob.MimeType != "" && blob.MimeType != "unknown" {
			w.Header().Add("Content-Type", blob.MimeType)
		}
		w.Header().Add("Content-Disposition", fmt.Sprintf("%s;filename=%s", "attachment", blob.FileName))
		rdr := bytes.NewReader(blob.File)
		http.ServeContent(w, r, blob.FileName, time.Now(), rdr)
		return
	}
	return fn
}
func (c *API) setPasswordHandler() func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		type Request struct {
			Password string `json:"password"`
		}
		req := &Request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		u.PasswordHash = HashPassword(req.Password)
		_, err = u.UpdateG(boil.Whitelist(db.UserColumns.PasswordHash))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 200, nil
	}
	return fn
}

func (c *API) attendanceListHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	IntegrationIDStr := chi.URLParam(r, "integration_id")
	TeacherIDStr := chi.URLParam(r, "teacher_id")
	type Response struct {
		Data db.AttendanceSlice `json:"data"`
	}

	IntegrationID, err := strconv.Atoi(IntegrationIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = isIntegrationOwner(IntegrationID, u.ID)
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	TeacherID, err := strconv.Atoi(TeacherIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	result, err := db.Attendances(
		db.AttendanceWhere.IntegrationID.EQ(null.Int64From(int64(IntegrationID))),
		db.AttendanceWhere.TeacherID.EQ(null.Int64From(int64(TeacherID))),
	).AllG()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &Response{result}, 200, nil
}
func (c *API) friendListHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	IntegrationIDStr := chi.URLParam(r, "integration_id")
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	IntegrationID, err := strconv.Atoi(IntegrationIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = isIntegrationOwner(IntegrationID, u.ID)
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	// err = refreshFriendCache(IntegrationID)
	// if err != nil {
	// 	return nil, http.StatusInternalServerError, err
	// }
	result, err := db.Friends(db.FriendWhere.IntegrationID.EQ(int64(IntegrationID))).AllG()
	if err != nil {
		return nil, 500, err
	}
	if result == nil {
		return &Response{}, 500, err
	}
	return &Response{result}, 200, nil
}
func (c *API) friendRefreshHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	type Response struct {
		Data db.FriendSlice `json:"data"`
	}
	IntegrationIDStr := chi.URLParam(r, "integration_id")
	IntegrationID, err := strconv.Atoi(IntegrationIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = isIntegrationOwner(IntegrationID, u.ID)
	if err != nil {
		return nil, http.StatusForbidden, err
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
func (c *API) friendPromoteHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	type Response struct {
		Data *db.Friend `json:"data"`
	}

	IntegrationIDStr := chi.URLParam(r, "integration_id")
	IntegrationID, err := strconv.Atoi(IntegrationIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = isIntegrationOwner(IntegrationID, u.ID)
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	FriendID := chi.URLParam(r, "friend_id")
	friend, err := db.Friends(
		db.FriendWhere.IntegrationID.EQ(int64(IntegrationID)),
		db.FriendWhere.VrchatID.EQ(FriendID),
	).OneG()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	friend.IsTeacher = true
	_, err = friend.UpdateG(boil.Whitelist(db.FriendColumns.IsTeacher))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return friend, 200, nil
}
func (c *API) friendDemoteHandler(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	type Response struct {
		Data *db.Friend `json:"data"`
	}

	IntegrationIDStr := chi.URLParam(r, "integration_id")
	IntegrationID, err := strconv.Atoi(IntegrationIDStr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = isIntegrationOwner(IntegrationID, u.ID)
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	FriendID := chi.URLParam(r, "friend_id")
	friend, err := db.Friends(
		db.FriendWhere.IntegrationID.EQ(int64(IntegrationID)),
		db.FriendWhere.VrchatID.EQ(FriendID),
	).OneG()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	friend.IsTeacher = false
	_, err = friend.UpdateG(boil.Whitelist(db.FriendColumns.IsTeacher))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return friend, 200, nil
}
func (c *API) teacherListHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
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

func (c *API) metricsHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, 200, ErrNotImplemented
}

func (c *API) userJWTHandler(auther *Auther) func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		type Response struct {
			Data string `json:"data"`
		}
		expiration := time.Now().Add(time.Duration(30) * time.Hour * 24)
		jwt, err := auther.GenerateJWT(u.Email, strconv.Itoa(int(u.ID.Int64)), u.Role, expiration)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return &Response{jwt}, http.StatusOK, nil
	}
	return fn
}
func (c *API) userListHandler() func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		if u.Role != roleAdmin {
			return nil, http.StatusForbidden, errors.New("unauthorized")
		}
		type Response struct {
			Data db.UserSlice `json:"data"`
		}
		users, err := db.Users().AllG()
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		return &Response{users}, 200, nil
	}
	return fn
}
func (c *API) apiKeyHandler(auther *Auther) func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		type Response struct {
			Token string `json:"token"`
		}
		expiration := time.Now().Add(time.Duration(30) * time.Hour * 24)
		id := strconv.Itoa(int(u.ID.Int64))
		token, err := auther.GenerateJWT(u.Email, id, u.Role, expiration)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		return &Response{token}, http.StatusOK, nil
	}
	return fn
}
func (c *API) userImpersonateHandler(auther *Auther) func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
	fn := func(w http.ResponseWriter, r *http.Request, u *db.User) (interface{}, int, error) {
		if u.Role != roleAdmin {
			return nil, http.StatusForbidden, errors.New("unauthorized")
		}
		type Response struct {
			Token string `json:"token"`
		}
		targetUserIDStr := chi.URLParam(r, "user_id")
		targetUserID, err := strconv.Atoi(targetUserIDStr)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		targetUser, err := db.FindUserG(null.Int64From(int64(targetUserID)))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		expiration := time.Now().Add(time.Duration(30) * time.Hour * 24)
		jwt, err := auther.GenerateJWT(targetUser.Email, strconv.Itoa(int(targetUser.ID.Int64)), targetUser.Role, expiration)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		cookie := http.Cookie{Name: "jwt", Value: jwt, Expires: expiration, HttpOnly: true, Path: "/", SameSite: http.SameSiteDefaultMode, Secure: false}
		http.SetCookie(w, &cookie)

		return &Response{jwt}, http.StatusOK, nil
	}
	return fn
}
