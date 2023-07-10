package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"net/http"
	"strings"
	"time"

	upper "github.com/upper/db/v4"
)

type Token struct {
	ID        int    `db:"id,omitempty" json:"id"`
	UserID    int    `db:"user_id" json:"user_id"`
	FirstName string `db:"first_name" json:"first_name"`
	Email     string `db:"email" json:"email"`
	PlainText string `db:"token" json:"token"`
	Hash      string `db:"token_hash" json:"-"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	Expires   string `db:"expiry" json:"expiry"`
}

func (t *Token) Table() string {
	return "tokens"
}

func (t *Token) GenerateToken(userID int, ttl time.Duration) (token *Token, err error) {
	token = &Token{
		UserID:  userID,
		Expires: time.Now().Add(ttl).Format(time.RFC3339),
	}

	randomBytes := make([]byte, 16)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))

	token.Hash = string(hash[:])

	return token, nil
}

func (t *Token) InsertToken(token Token, u User) error {
	collection := database.Collection(t.Table())

	// delete existing tokens
	res := collection.Find(upper.Cond{"user_id =": u.ID})
	err := res.Delete()
	if err != nil {
		return err
	}

	// add new tokens
	token.CreatedAt = time.Now().Format(time.RFC3339)
	token.UpdatedAt = time.Now().Format(time.RFC3339)
	token.FirstName = u.FirstName
	token.Email = u.Email

	_, err = collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) GetUserForToken(token string) (u *User, err error) {
	var theToken Token

	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"token =": token})
	err = res.One(&theToken)
	if err != nil {
		return nil, err
	}

	collection = database.Collection(u.Table())
	res = collection.Find(upper.Cond{"id =": theToken.UserID})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}

	u.Token = theToken

	return u, nil
}

func (t *Token) GetTokensForUser(id int) (tokens []*Token, err error) {
	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"user_id =": id})
	err = res.All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (t *Token) GetTokenByID(id int) (token *Token, err error) {
	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"id =": id})
	err = res.One(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *Token) GetTokenByEmail(email string) (token *Token, err error) {
	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"email =": email})
	err = res.One(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *Token) GetTokenByToken(hash string) (token *Token, err error) {
	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"token_hash =": hash})
	err = res.One(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *Token) AuthenticateToken(r *http.Request) (u *User, err error) {
	// get authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	// extract authorization token
	authToken := headerParts[1]

	if len(authToken) != 26 {
		return nil, errors.New("wrong token")
	}

	hash := sha256.Sum256([]byte(authToken))

	// check if token exists
	token, err := t.GetTokenByToken(string(hash[:]))
	if err != nil {
		return nil, errors.New("wrong token")
	}

	// check token expiry
	expiry, err := time.Parse(time.RFC3339, token.Expires)
	if err != nil {
		return nil, errors.New("error parsing time from database")
	}
	if expiry.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	// get user of token
	u, err = t.GetUserForToken(authToken)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (t *Token) ValidateToken(token string) (bool, error) {
	// get user of token
	u, err := t.GetUserForToken(token)
	if err != nil {
		return false, errors.New("user not found")
	}

	// check for a token
	if u.Token.PlainText == "" {
		return false, errors.New("token not found")
	}
	// check token expiry
	expiry, err := time.Parse("2006-01-01T00:00:00.000", u.Token.Expires)
	if err != nil {
		return false, errors.New("error parsing time from database")
	}
	if expiry.Before(time.Now()) {
		return false, errors.New("expired token")
	}

	return true, nil
}

func (t *Token) DeleteTokenByID(id int) error {
	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"id =": id})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteTokenByToken(plainText string) error {
	collection := database.Collection(t.Table())
	res := collection.Find(upper.Cond{"token =": plainText})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}
