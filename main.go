package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	handlers "github.com/mcg25035/assignment/handlers"
	dataHandler "github.com/mcg25035/assignment/handlers/dataHandler"
)


var app = fiber.New(&fiber.Settings{
	Prefork: false,
})

func main() {
	var apiv1 = app.Group("/api/v1").(*fiber.Group)
	fmt.Println(app.IsChild())
	dataHandler.InitDataHandler()
	handlers.RegisterAPI(apiv1, app.IsChild())
	app.Listen(":3000")
	return
}
