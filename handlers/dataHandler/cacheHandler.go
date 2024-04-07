package dataHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	// "time"
	// "strings"
	"strconv"

	"github.com/mcg25035/assignment/conditions"
	"github.com/mcg25035/assignment/utils"
	"github.com/redis/go-redis/v9"
    // "github.com/valyala/fastjson"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

func cacheGenerateID(baseKey string) string {
	id, _ := rdb.Incr(ctx, baseKey+":counter").Result()
	specialId := fmt.Sprintf("%s:%d", baseKey, id)
	return specialId
}

func cacheAddAge(age int, conditionID int64) {
    var key = "age:"+strconv.Itoa(age)
    var err = rdb.SAdd(ctx, key, conditionID).Err()
    if err != nil {
        panic(err)
    }
}

func cacheUpdateAge(age int, conditions []int64) {
    var key = "age:"+strconv.Itoa(age)
    for _, condition := range conditions {
        var err = rdb.SAdd(ctx, key, condition).Err()
        if err != nil {
            panic(err)
        }
    }
}

func cacheAddGender(gender conditions.Gender, conditionID int64) {
    var key = "gender:"+string(gender)
    var err = rdb.SAdd(ctx, key, conditionID).Err()
    if err != nil {
        panic(err)
    }
}

func cacheUpdateCountry(country conditions.Country, conditions []int64) {
    var key = "country:"+string(country)
    for _, condition := range conditions {
        var err = rdb.SAdd(ctx, key, condition).Err()
        if err != nil {
            panic(err)
        }
    }
}

func cacheAddCountry(country conditions.Country, conditionID int64) {
    var key = "country:"+string(country)
    var err = rdb.SAdd(ctx, key, conditionID).Err()
    if err != nil {
        panic(err)
    }
}

func cacheUpdateGender(gender conditions.Gender, conditions []int64) {
    var key = "gender:"+string(gender)
    for _, condition := range conditions {
        var err = rdb.SAdd(ctx, key, condition).Err()
        if err != nil {
            panic(err)
        }
    }
}

func cacheAddPlatform(platform conditions.Platform, conditionID int64) {
    var key = "platform:"+string(platform)
    var err = rdb.SAdd(ctx, key, conditionID).Err()
    if err != nil {
        panic(err)
    }
}

func cacheUpdatePlatform(platform conditions.Platform, conditions []int64) {
    var key = "platform:"+string(platform)
    for _, condition := range conditions {
        var err = rdb.SAdd(ctx, key, condition).Err()
        if err != nil {
            panic(err)
        }
    }
}


func cacheAddAd(adID string, ad AdData, expiration time.Time) {
    var key = "ad_data:"+adID
    var adData = "{\"title\":\""+ad.Title+"\",\"endAt\":\""+ad.EndAt+"\",\"startAt\":\""+ad.StartAt+"\"}"
    rdb.Set(ctx, key, adData, 0)
    rdb.ExpireAt(ctx, key, expiration)
}

func cacheAddCondition(adID string, conditionID int64, expiration time.Time) {
    var key = "conditions:"+strconv.FormatInt(conditionID, 10)
    rdb.Set(ctx, key, adID, 0)
    rdb.ExpireAt(ctx, key, expiration)
}

func cacheGetAd(
    age int,
    gender conditions.Gender,
    country conditions.Country,
    platform conditions.Platform,
    offset int,
    limit int,
) []AdData {
    var timeNow = time.Now().UnixMilli()
    var conditionsResult = rdb.SInter(ctx,
        "age:"+strconv.Itoa(age),
        "gender:"+string(gender),
        "country:"+string(country),
        "platform:"+string(platform),
    ).Val()
    

    for i := 0; i < len(conditionsResult); i++ {
        conditionsResult[i] = "conditions:"+conditionsResult[i]
    }
    
    var adResult = rdb.MGet(ctx, conditionsResult...).Val()
    var mapResult = make(map[string]bool)
    var keys = make([]string, 0, len(adResult))
    for _, adID := range adResult {
        if adID == nil {continue}
        _, exist := mapResult[adID.(string)]
        if exist {continue}
        mapResult[adID.(string)] = true
        keys = append(keys, "ad_data:"+adID.(string))
    }

    var adDataResult = rdb.MGet(ctx, keys...).Val()

    
    var result = []AdData{};
    for _, ad := range adDataResult {
        if ad == nil {continue}
        var redisResponse AdDataWithStartAt
        var adData AdData
        json.Unmarshal([]byte(ad.(string)), &redisResponse)
        var adDataStartAt, _ = utils.DateStringToTimestamp(redisResponse.StartAt)
        if adDataStartAt > timeNow {continue}
        adData.Title = redisResponse.Title
        adData.EndAt = redisResponse.EndAt
        result = append(result, adData)
    }

    sort.Slice(result, func(i, j int) bool {
        return result[i].Title < result[j].Title
    })

    if offset > len(result) {
        return []AdData{}
    }

    if (offset > 0) {
        result = result[offset:]
    }

    if limit < len(result) {
        result = result[:limit]
    } 

    return result
}
