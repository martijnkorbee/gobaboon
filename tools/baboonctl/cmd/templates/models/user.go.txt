package models

import (
	"errors"
	"time"

	upper "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `db:"id,omitempty"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Active    int    `db:"user_active"`
	Password  string `db:"password"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	Token     Token  `db:"-"`
}

func (u *User) Table() string {
	return "users"
}

// AddUser adds a user
func (u *User) AddUser(user User) (ID int, err error) {
	// check if email is available
	usr, _ := u.GetUserByEmail(user.Email)
	if usr != nil {
		return 0, errors.New("email address not available")
	}

	// create new password hash
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.CreatedAt = time.Now().Format(time.RFC3339)
	user.UpdatedAt = time.Now().Format(time.RFC3339)
	user.Password = string(newHash)

	collection := database.Collection(u.Table())
	res, err := collection.Insert(user)
	if err != nil {
		return 0, err
	}
	ID = GetInsertID(res.ID())

	return ID, nil
}

// GetAll gets all users
func (u *User) GetAll() ([]*User, error) {
	var all []*User
	collection := database.Collection(u.Table())

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetUserByID gets user by ID also adds the token, token validation should still be performed
func (u *User) GetUserByID(ID int) (*User, error) {
	var user User
	collection := database.Collection(u.Table())

	res := collection.Find(upper.Cond{"id =": ID})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = database.Collection(token.Table())

	res = collection.Find(upper.Cond{"user_id =": user.ID}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != upper.ErrNilRecord && err != upper.ErrNoMoreRows {
			return nil, err
		}
	}

	user.Token = token

	return &user, nil
}

// GetUserByEmail gets user by email also adds the token, token validation should still be performed
func (u *User) GetUserByEmail(email string) (*User, error) {
	var user User
	collection := database.Collection(u.Table())

	res := collection.Find(upper.Cond{"email =": email})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	// add token to user
	var token Token
	collection = database.Collection(token.Table())

	res = collection.Find(upper.Cond{"user_id =": user.ID}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != upper.ErrNilRecord && err != upper.ErrNoMoreRows {
			return nil, err
		}
	}

	user.Token = token

	return &user, nil
}

// UpdateUser updates a user
func (u *User) UpdateUser(user User) error {
	user.UpdatedAt = time.Now().String()

	collection := database.Collection(u.Table())
	res := collection.Find(user.ID)
	err := res.Update(&user)
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword resets the users password
func (u *User) ResetPassword(id int, password string) error {
	// create new password hash
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Password = string(newHash)

	err = user.UpdateUser(*user)
	if err != nil {
		return err
	}

	return nil
}

// AuthenticateUser authenticates user
func (u *User) AuthenticateUser(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// DeleteUser deletes a user
func (u *User) DeleteUser(id int) error {
	var user User

	collection := database.Collection(u.Table())
	res := collection.Find(id)

	// check if user exists
	err := res.One(&user)
	if err != nil {
		return err
	}

	// delete user
	err = res.Delete()
	if err != nil {
		return err
	}

	return nil
}