package main

func main() {
	app := newApp()
	if err := app.run(); err != nil {
		panic(err)
	}
}
