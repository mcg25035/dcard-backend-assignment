package handlers

import (
	_ "fmt"
	"sync"

	"github.com/mcg25035/assignment/conditions"
	"github.com/gofiber/fiber"
	"github.com/mcg25035/assignment/handlers/dataHandler"
)

var mutex sync.Mutex

func RegisterAPI(api *fiber.Group, isChild bool) {
	api.Get(("/ad"), func(c *fiber.Ctx) {
		type AdRequest struct {
			Age		int						`query:"age"`
			Offset  int						`query:"offset"`
			Limit   int						`query:"limit"`
			Gender  conditions.Gender		`query:"gender"`
			Country conditions.Country		`query:"country"`
			Platform conditions.Platform	`query:"platform"`
		}

		var request AdRequest
		if err := c.QueryParser(&request); err != nil {
			c.Status(400).Send("Invalid query")
			return
		}

		var limit = request.Limit
		var offset = request.Offset
		var age = request.Age
		var gender = request.Gender
		var country = request.Country
		var platform = request.Platform

		var data = dataHandler.GetAdData(age, gender, country, platform, limit, offset)
		// fmt.Println(offset)
		// fmt.Println(data)
		c.JSON(data)
	})

	api.Get("/test", func(c *fiber.Ctx) {
		c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	if (isChild) {
		return
	}

	api.Post(("/ad"), func(c *fiber.Ctx) {
		dataHandler.DbLockOpreate()
		defer dataHandler.DbUnlockOpreate()

		var request dataHandler.AdRequest
		if err := c.BodyParser(&request); err != nil {
			c.Status(400).Send("Invalid JSON")
			return
		}

		dataHandler.InsertAd(request)
		c.JSON(fiber.Map{
			"status": "OK",
		})
	})
}
