package main

import (
	"bless-activity/application"
)

func main() {
	app := application.NewApp()

	if err := app.Start(); err != nil {
		panic(err)
	}
}
