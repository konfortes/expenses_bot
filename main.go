package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// auto load .env file
	_ "github.com/joho/godotenv/autoload"
)

type App struct {
	DB string
}

func (app *App) init() {
	fmt.Println("initializing app...")
}

var app *App
var weatherClient WeatherClient

func main() {
	app = &App{}
	app.init()
	weatherClient = WeatherClient{URL: os.Getenv("WEATHER_API_URL"), Key: os.Getenv("WEATHER_API_KEY")}
	fmt.Println("listening on", os.Getenv("PORT"))
	go http.ListenAndServe(fmt.Sprintf("localhost:%s", os.Getenv("PORT")), nil)
	err := initBot()
	if err != nil {
		log.Fatal(err)
		panic("error initializing bot")
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Fatal("panic: ", r)
	}
}
