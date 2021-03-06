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
		log.Fatalln("error initializing application: ", err)
	}

	weatherClient = WeatherClient{URL: os.Getenv("WEATHER_API_URL"), Key: os.Getenv("WEATHER_API_KEY")}

	listenURL := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	fmt.Println("listening on", listenURL)
	// TODO: understand why go?
	go http.ListenAndServe(listenURL, nil)

	err := initBot()
	if err != nil {
		log.Fatalln("error initializing bot: ", err)
	}

}

func handlePanic() {
	if r := recover(); r != nil {
		log.Print("panic: ", r)
	}
}
