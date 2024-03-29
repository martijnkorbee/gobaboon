package models

import (
	"fmt"
	"github.com/martijnkorbee/gobaboon/pkg/db"
	upper "github.com/upper/db/v4"
)

var (
	database upper.Session
)

type Models struct {
	// any models specified here (and added in the New function)
	// are accessible to the application
	/*
		<- required after you use baboonctl make auth ->
		//Users  User
		//Tokens Token
	*/
}

func New(db *db.Database) *Models {
	database = db.Session

	return &Models{
		/*
			<- required after you use baboonctl make auth ->
			//Users:  User{},
			//Tokens: Token{},
		*/
	}
}

func GetInsertID(i upper.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
