package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// auto load .env file
	_ "github.com/joho/godotenv/autoload"
)

type App struct {
	DB *sql.DB
}

func (app *App) init() error {
	var err error
	app.DB, err = initDB()
	return err
}

var app *App
var weatherClient WeatherClient

func main() {
	app = &App{}
	if err := app.init(); err != nil {
		panic(err)
	}

	weatherClient = WeatherClient{URL: os.Getenv("WEATHER_API_URL"), Key: os.Getenv("WEATHER_API_KEY")}
	port := os.Getenv("PORT")
	fmt.Println("listening on", port)
	go http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil)
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
