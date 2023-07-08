package main

import (
	"github.com/martijnkorbee/gobaboon/internal/app"
)

func main() {
	application := app.New()
	application.Start()
}
