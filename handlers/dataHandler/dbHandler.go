package dataHandler

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mcg25035/assignment/conditions"
	"github.com/mcg25035/assignment/utils"
	"sync"
)

var db, err = sql.Open("sqlite3", "file:./ad_data.db?cache=shared")
var dbOpreatingLock = sync.RWMutex{}

func DbLockOpreate() {
	dbOpreatingLock.Lock()
}

func DbUnlockOpreate() {
	dbOpreatingLock.Unlock()
}

func dbCreateTable(table_name string, columns string){
	var sqlCmd = `create table if not exists ` + table_name + ` (` + columns + `);`
	_, err = db.Exec(sqlCmd)
	if err != nil {
		panic(err)
	}
}

func dbGetID() string {
	var id = utils.GenerateRandomString()
	for {
		var rows, err = db.Query("select * from ad_data where id = ?", id)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		if rows.Next() {
			id = utils.GenerateRandomString()
		} else {
			break
		}
	}
	return id
}

func dbInit() {
	if err != nil {
		panic(err)
	}

	dbCreateTable("ad_data", `
		id text,
		title text,
		startAt int,
		endAt int
	`)

	dbCreateTable("ad_condition", `
		ref_id text,
		ageStart int,	
		ageEnd int,
		gender text,
		country text,
		platform text,
		startAt int,
		endAt int
	`)

	// db.SetMaxOpenConns(1)

	utils.Log("Database initialized.")
}

func dbGetFuture5min() int64 {
	return time.Now().Add(5 * time.Minute).UnixMilli()
}

func adAddAD (title string, startAt int64, endAt int64) string {
	var id = dbGetID()
	var sqlCmd = `insert into ad_data values (?, ?, ?, ?);`
	db.Exec(
		sqlCmd,
		id,
		title,
		startAt,
		endAt,
	)
	return id
}

func dbAddCondition (
	id string,
	ageStart int, 
	ageEnd int, 
	gender []conditions.Gender,
	country []conditions.Country,
	platform []conditions.Platform,
) int64 {
	var sqlCmd = `insert into ad_condition values (?, ?, ?, ?, ?, ?, ?, ?);`
	var adData = dbGetAdById(id)
	var startAt, _ = utils.DateStringToTimestamp(adData.StartAt)
	var endAt, _ = utils.DateStringToTimestamp(adData.EndAt)

	var result, err = db.Exec(
		sqlCmd,
		id,
		ageStart,
		ageEnd,
		utils.GenderArrayToString(gender),
		utils.CountryArrayToString(country),
		utils.PlatformArrayToString(platform),
		startAt,
		endAt,
	)

	if err != nil {
		panic(err)
	}

	var lastInsertId, err0 = result.LastInsertId()

	if err0 != nil {
		panic(err)
	}

	return lastInsertId
}

func dbAgeFilter (age int) []int64 {
	var result = []int64{}
	var future = dbGetFuture5min()
	var now = getNowTimestamp()
	var sqlCmd = `select rowid from ad_condition where ageStart <= ? and ageEnd >= ? 
		and ? <= endAt and startAt <= ?;
	`
	var rows, err = db.Query(sqlCmd, age, age, now-1000, future)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var condition_id int64
		rows.Scan(&condition_id)
		result = append(result, int64(condition_id))
	}
	return result
}

func dbGenderFilter (gender conditions.Gender) []int64 {
	var result = []int64{}
	var now = getNowTimestamp()
	var future = dbGetFuture5min()
	var sqlCmd = `select rowid from ad_condition where gender like ?
		and ? <= endAt and startAt <= ?;
	`
	var rows, err = db.Query(sqlCmd, "%"+ string(gender) +"%", now-1000, future)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var condition_id int64
		rows.Scan(&condition_id)
		result = append(result, int64(condition_id))
	}
	return result
}

func dbCountryFilter (country conditions.Country) []int64 {
	var result = []int64{}
	var now = getNowTimestamp()
	var future = dbGetFuture5min()
	var sqlCmd = `select rowid from ad_condition where country like ?
		and ? <= endAt and startAt <= ?;
	`
	var rows, err = db.Query(sqlCmd, "%"+ string(country) +"%", now-1000, future)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var condition_id int64
		rows.Scan(&condition_id)
		result = append(result, condition_id)
	}
	return result
}

func dbPlatformFilter (platform conditions.Platform) []int64 {
	var result = []int64{}
	var now = getNowTimestamp()
	var future = dbGetFuture5min()
	var sqlCmd = `select rowid from ad_condition where platform like ?
		and ? <= endAt and startAt <= ?;
	`
	var rows, err = db.Query(sqlCmd, "%"+ string(platform) +"%", now-1000, future)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var condition_id int64
		rows.Scan(&condition_id)
		result = append(result, condition_id)
	}
	return result
}

func dbGetConditionAdMap () (map[int64]string, []int64) {
	var result = map[int64]string{}
	var resultExpire = []int64{}
	var now = getNowTimestamp()
	var future = dbGetFuture5min()
	var sqlCmd = `select rowid, ref_id, endAt from ad_condition where ? <= endAt and startAt <= ?;`
	var rows, err = db.Query(sqlCmd, now-1000, future)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var condition_id int64
		var ad_id string
		var expireAt int64
		rows.Scan(&condition_id, &ad_id, &expireAt)
		result[condition_id] = ad_id
		resultExpire = append(resultExpire, expireAt)
	}
	return result, resultExpire
}

func dbGetAdDataIdMap() map[string]AdData {
	var result = map[string]AdData{}
	var now = getNowTimestamp()
	var future = dbGetFuture5min()
	var sqlCmd = `select id, title, startAt, endAt from ad_data 
		where ? <= endAt and startAt <= ?;
	`
	var rows, err = db.Query(sqlCmd, now-1000, future)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ad_id string
		var title string
		var startAt int64
		var endAt int64
		rows.Scan(&ad_id, &title, &startAt, &endAt)
		result[ad_id] = AdData{
			title, 
			utils.TimestampToDateString(endAt), 
			utils.TimestampToDateString(startAt),
		}
	}
	return result
}


func dbGetAdById (id string) AdData {
	var sqlCmd = `select title, startAt, endAt from ad_data where id = ?;`
	var row, err = db.Query(sqlCmd, id)
	if err != nil {
		panic(err)
	}
	if (!row.Next()) {
		return AdData{}
	}
	var title string
	var startAt int64
	var endAt int64
	row.Scan(&title, &startAt, &endAt)
	return AdData{title, utils.TimestampToDateString(endAt), utils.TimestampToDateString(startAt)}
}

// func getAdData_ (
// 	age int,
// 	gender conditions.Gender,
// 	country conditions.Country, 
// 	platform conditions.Platform,
// 	limit int,
// 	offset int,
// ) []AdData {
// 	var result = []AdData{};
// 	var idResult = []string{};
// 	var sqlCmd = `select ref_id from ad_condition where
// 		ageStart <= ? and 
// 		ageEnd >= ? and 
// 		gender like ? and
// 		country like ? and
// 		platform like ?;`
	
// 	var rows, err = db.Query(
// 		sqlCmd,
// 		age,
// 		age,
// 		"%"+string(gender)+"%",
// 		"%"+string(country)+"%",
// 		"%"+string(platform)+"%",
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	var offsetCounter = 0
// 	var limitCounter = 0

// 	var now = getNowTimestamp()

// 	defer rows.Close();
// 	for rows.Next(){
// 		var id string
// 		rows.Scan(&id)
// 		if (utils.Contains(idResult, id)) {continue}
// 		idResult = append(idResult, id)
// 		if (offsetCounter < offset) {
// 			offsetCounter++
// 			continue
// 		}
// 		limitCounter++
// 		if (limitCounter > limit) {break}
// 		var row, err = db.Query("select title, startAt, endAt from ad_data where id = ?", id)
// 		if err != nil {panic(err)}
// 		if (!row.Next()) {continue}
// 		var title string
// 		var startAt int64
// 		var endAt int64
// 		row.Scan(&title, &startAt, &endAt)
// 		if (startAt > now) {continue}
// 		if (endAt < now) {continue}			
// 		result = append(result, AdData{title, utils.TimestampToDateString(endAt)})
// 	}

// 	return result
// }

func getNowTimestamp() int64 {
	return time.Now().UnixMilli()	
}
