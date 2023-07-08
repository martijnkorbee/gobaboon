package main

import (
	"github.com/martijnkorbee/gobaboon/cmd/app"
)

func main() {
	application := app.MustInitApplication()
	application.Start()
}
