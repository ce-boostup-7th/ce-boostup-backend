package main

import (
	"ce-boostup-backend/db"
	"ce-boostup-backend/route"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// loads values from  .env into the system
	if err := godotenv.Load("./variables.env"); err != nil {
		log.Fatal("No .env file found")
	}

	e := route.Init()
	db.Init()
	e.Logger.Fatal(e.Start(":1323"))
}
