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
	// any models inserted here (and must be inserted through the New function)
	// are easily accessible throughout the entire application
}

func New(db *db.Database) *Models {
	database = db.Session

	return &Models{}
}

func getInsertID(i upper.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
