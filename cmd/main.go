package main

import (
	"baboonapp/cmd/web/app"
)

func main() {
	app := app.MustInitApplication()
	app.Start()
}
