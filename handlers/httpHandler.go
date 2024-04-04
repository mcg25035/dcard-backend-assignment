package handlers

import (
	"github.com/gofiber/fiber"
	"github.com/mcg25035/assignment/conditions"
	"github.com/mcg25035/assignment/utils"
)

func RegisterAPI(api *fiber.Group) {
	api.Post(("/ad"), func(c *fiber.Ctx) {
		type AdRequest struct {
			Title      string    `json:"title"`
			StartAt    string    `json:"startAt"` 
			EndAt      string    `json:"endAt"`   
			Conditions[] struct {
				AgeStart  int      `json:"ageStart"`
				AgeEnd    int      `json:"ageEnd"`
				Country   []conditions.Country `json:"country"`
				Platform  []conditions.Platform `json:"platform"`
				Gender  conditions.Gender `json:"gender"`
			} `json:"conditions"`
		}

		var request AdRequest
		if err := c.BodyParser(&request); err != nil {
			c.Status(400).Send("Invalid JSON")
			return
		}

		var formattedStartAt, err1 = utils.DateStringToTimestamp(request.StartAt)
		var formattedEndAt, err2 = utils.DateStringToTimestamp(request.EndAt)		

		if (err1 != nil || err2 != nil) {
			c.Status(400).Send("Invalid date format")
			return
		}

		var adId = InsertAdData(request.Title, formattedStartAt, formattedEndAt)
		if (len(request.Conditions) == 0){
			InsertAdCondition(
				adId,
				1,
				100,
				conditions.AllGenders(),
				conditions.AllCountries(),
				conditions.AllPlatforms(),
			)
		}
		for _, condition := range request.Conditions {
			var ageStart = condition.AgeStart
			var ageEnd = condition.AgeEnd
			var gender = make([]conditions.Gender, 0)
			if condition.Gender == "" {
				gender = conditions.AllGenders()
			} else {
				gender = append(gender, condition.Gender)
			}
			var country = condition.Country
			var platform = condition.Platform

			if (ageStart == 0 && ageEnd == 0){
				ageStart = 1
				ageEnd = 100
			}

			// if (len(gender) == 0){
			// 	gender = conditions.AllGenders()
			// }

			if (len(country) == 0){
				country = conditions.AllCountries()
			}

			if (len(platform) == 0){
				platform = conditions.AllPlatforms()
			}

			InsertAdCondition(
				adId, 
				ageStart,
				ageEnd, 
				gender, 
				country, 
				platform,
			)

			c.JSON(fiber.Map{
				"status": "OK",
			})
		}
	})

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

		var data = GetAdData(age, gender, country, platform, limit, offset)
		c.JSON(data)


	})
}
