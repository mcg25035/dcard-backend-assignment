package main

import (
	"github.com/gofiber/fiber"
	handlers "github.com/mcg25035/assignment/handlers"
)

var app = fiber.New()

func main() {
	var apiv1 = app.Group("/api/v1").(*fiber.Group)
	handlers.RegisterAPI(apiv1)
	handlers.InitDatabase()	
	app.Listen(":3000")
}
