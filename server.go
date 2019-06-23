package main

import (
	"ce-boostup-backend/db"
	"ce-boostup-backend/route"
)

func main() {
	e := route.Init()
	db.Init()
	e.Logger.Fatal(e.Start(":1323"))
}
