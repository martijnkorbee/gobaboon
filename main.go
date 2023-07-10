package main

import (
	"github.com/martijnkorbee/gobaboon/internal/application"
)

var (
	app *application.Application
)

func init() {
	app = application.New()
}

func main() {
	app.Run()
}
