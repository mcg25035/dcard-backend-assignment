package dataHandler

import (
	// "strconv"

	_ "fmt"
	"time"

	"github.com/mcg25035/assignment/conditions"
	"github.com/mcg25035/assignment/utils"
)

func intervalSynchonizer(){
	dbOpreatingLock.Lock()
	defer dbOpreatingLock.Unlock()
	for i := 1; i < 101; i++ {
		var conditions = dbAgeFilter(i)
		cacheUpdateAge(i, conditions)
	}
	for _, g := range conditions.AllGenders() {
		var conditions = dbGenderFilter(g)
		cacheUpdateGender(g, conditions)
	}
	for _, c := range conditions.AllCountries() {
		var conditions = dbCountryFilter(c)
		cacheUpdateCountry(c, conditions)
	}
	for _, p := range conditions.AllPlatforms() {
		var conditions = dbPlatformFilter(p)
		cacheUpdatePlatform(p, conditions)
	}
	var i = 0
	var conditionAdMap, exipreAt = dbGetConditionAdMap()
	for conditionID, adID := range conditionAdMap {
		cacheAddCondition(adID, conditionID, utils.TimestampToTime(exipreAt[i]))
		i++
	}
	i = 0
	for adID, adData := range dbGetAdDataIdMap() {
		cacheAddAd(adID, adData, utils.TimestampToTime(exipreAt[i]))
		i++
	}
	
	utils.Log("Synchonized cache with database.")
}

func RegisterSynchonizer() {
	ticker := time.NewTicker(3 * time.Minute)
	go func() {
		for range ticker.C {
			intervalSynchonizer()
		}
	}()
}

func InitDataHandler() {
	dbInit()
	rdb.FlushAll(ctx)
	intervalSynchonizer()
	utils.Log("Redis cache initialized.")
	RegisterSynchonizer()
}

func addCondition(adCondition AdCondition, conditionID int64){
	var ageStart = adCondition.AgeStart
	var ageEnd = adCondition.AgeEnd
	var gender = adCondition.Gender
	var country = adCondition.Country
	var platform = adCondition.Platform
	
	for i := ageStart; i <= ageEnd; i++ {
		cacheAddAge(i, conditionID)
	}

	for _, g := range gender {
		cacheAddGender(g, conditionID)
	}
	
	for _, c := range country {
		cacheAddCountry(c, conditionID)
	}

	for _, p := range platform {
		cacheAddPlatform(p, conditionID)
	}
}

func addAd(title string, startAt int64, endAt int64, adConditions []AdCondition) {
	var id = adAddAD(title, startAt, endAt)
	var expiration = utils.TimestampToTime(endAt)
	for _, condition := range adConditions {
		var conditionID = dbAddCondition(
			id, 
			condition.AgeStart, 
			condition.AgeEnd, 
			condition.Gender,
			condition.Country,
			condition.Platform,
		)
		addCondition(condition, conditionID)
		cacheAddCondition(id, conditionID, expiration)
	}

	cacheAddAd(id, AdData{
		Title: title,
		StartAt: utils.TimestampToDateString(startAt),
		EndAt: utils.TimestampToDateString(endAt),
	}, expiration)
}

func formatConditions(conditionsRaw []AdConditionRaw) []AdCondition {
	var result = []AdCondition{}

	for _, condition := range conditionsRaw {
		var formatedCondition = AdCondition{}
		if (condition.AgeStart == 0 && condition.AgeEnd == 0){
			condition.AgeStart = 1
			condition.AgeEnd = 100
		}
		formatedCondition.AgeStart = condition.AgeStart
		formatedCondition.AgeEnd = condition.AgeEnd
		var gender = make([]conditions.Gender, 0)
		if condition.Gender == "" {
			gender = conditions.AllGenders()
		} else {
			gender = append(gender, condition.Gender)
		}
		formatedCondition.Gender = gender
		formatedCondition.Country = condition.Country
		formatedCondition.Platform = condition.Platform
		if (len(formatedCondition.Country) == 0){
			formatedCondition.Country = conditions.AllCountries()
		}
		if (len(formatedCondition.Platform) == 0){
			formatedCondition.Platform = conditions.AllPlatforms()
		}

		result = append(result, formatedCondition)
	}

	if len(result) == 0 {
		var generalCondition = AdCondition{
			AgeStart: 1,
			AgeEnd: 100,
			Gender: conditions.AllGenders(),
			Country: conditions.AllCountries(),
			Platform: conditions.AllPlatforms(),
		}
		result = append(result, generalCondition)
	}		
	
	return result
}

func formatRequest(adRequest AdRequest) (AdData, []AdCondition) {
	var adData AdData
	var adConditions = formatConditions(adRequest.Conditions)
	adData.Title = adRequest.Title
	adData.EndAt = adRequest.EndAt
	adData.StartAt = adRequest.StartAt
	return adData, adConditions
}

func InsertAd(adRequest AdRequest) {
	var adData, adConditions = formatRequest(adRequest)
	
	var startAt, _ = utils.DateStringToTimestamp(adData.StartAt)
	var endAt, _ = utils.DateStringToTimestamp(adData.EndAt)

	addAd(adData.Title, startAt, endAt, adConditions)
}

func GetAdData (
	age int, 
	anti_policy_correctness conditions.Gender,
	country conditions.Country,
	platform conditions.Platform,
	limit int,
	offset int,
) []AdData {
	return cacheGetAd(age, anti_policy_correctness, country, platform, offset, limit)
}
