package accumulator

import (
	"accumulator/db"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/volatiletech/null"
)

// Auther to handle JWT authentication
type Auther struct {
	TokenAuth *jwtauth.JWTAuth
}

// NewAuther for JWT and blacklisting
func NewAuther(jwtsecret string) *Auther {
	result := &Auther{
		TokenAuth: jwtauth.New("HS256", []byte(jwtsecret), []byte(jwtsecret)),
	}
	return result
}

// FromContext grabs the user from the context if a JWT is inside
func (a *Auther) FromContext(ctx context.Context) (*db.User, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	idI, ok := claims["id"]
	if !ok {
		err := errors.New("could not get ID from claims")
		return nil, err
	}

	idStr := idI.(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	u, err := db.FindUserG(null.Int64From(int64(id)))
	if err != nil {
		return nil, err
	}

	return u, nil

}

// HashPassword encrypts a plaintext string and returns the hashed version in base64
func (a *Auther) HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(hashed)
}

// GenerateJWT returns the token for client side persistence
func (a *Auther) GenerateJWT(email, id string, expiration time.Time) (string, error) {
	_, tokenString, err := a.TokenAuth.Encode(jwt.MapClaims{"email": email, "id": id})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyMiddleware for authentication adds JWT to context down the HTTP chain
func (a *Auther) VerifyMiddleware() func(http.Handler) http.Handler {
	return jwtauth.Verifier(a.TokenAuth)
}

// ValidatePassword will check the login details
func (a *Auther) ValidatePassword(email string, password string) error {
	user, err := db.Users(db.UserWhere.Email.EQ(email)).OneG()
	if err != nil {
		return err
	}

	storedHash, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedHash, []byte(password))
	if err != nil {
		return err
	}

	return nil
}
