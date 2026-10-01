package main

import "github.com/xenobyter/xbf/internal/app"

func main() {
	a := app.New()
	if err := a.Run(); err != nil {
		panic(err)
	}
}
