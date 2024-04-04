package handlers

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mcg25035/assignment/conditions"
	"github.com/mcg25035/assignment/utils"
)

type AdData struct {
	Title string `json:"title"`
	EndAt string `json:"endAt"`
}

var db, err = sql.Open("sqlite3", "file:./ad_data.db?cache=shared")

func createTable(table_name string, columns string){
	var sqlCmd = `create table if not exists ` + table_name + ` (` + columns + `);`
	_, err = db.Exec(sqlCmd)
	if err != nil {
		panic(err)
	}
}

func getSafeID() string {
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

func InitDatabase() {
	if err != nil {
		panic(err)
	}

	createTable("ad_data", `
		id text,
		title text,
		startAt int,
		endAt int
	`)

	createTable("ad_condition", `
		ref_id text,
		ageStart int,	
		ageEnd int,
		gender text,
		country text,
		platform text
	`)

	fmt.Println("Database initialized")
}

func InsertAdData (title string, startAt int64, endAt int64) string {
	var id = getSafeID()
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

func InsertAdCondition (
	id string,
	ageStart int, 
	ageEnd int, 
	gender []conditions.Gender,
	country []conditions.Country,
	platform []conditions.Platform,
) {
	var sqlCmd = `insert into ad_condition values (?, ?, ?, ?, ?, ?);`

	var _, err = db.Exec(
		sqlCmd,
		id,
		ageStart,
		ageEnd,
		utils.GenderArrayToString(gender),
		utils.CountryArrayToString(country),
		utils.PlatformArrayToString(platform),
	)

	if err != nil {
		panic(err)
	}

	

}

func GetAdData (
	age int,
	gender conditions.Gender,
	country conditions.Country, 
	platform conditions.Platform,
	limit int,
	offset int,
) []AdData {
	var result = []AdData{};
	var idResult = []string{};
	var sqlCmd = `select ref_id from ad_condition where
		ageStart <= ? and 
		ageEnd >= ? and 
		gender like ? and
		country like ? and
		platform like ?;`
	
	var rows, err = db.Query(
		sqlCmd,
		age,
		age,
		"%"+string(gender)+"%",
		"%"+string(country)+"%",
		"%"+string(platform)+"%",
	)

	if err != nil {
		panic(err)
	}

	var offsetCounter = 0
	var limitCounter = 0

	var now = getNowTimestamp()

	defer rows.Close();
	for rows.Next(){
		var id string
		rows.Scan(&id)
		if (utils.Contains(idResult, id)) {continue}
		idResult = append(idResult, id)
		if (offsetCounter < offset) {
			offsetCounter++
			continue
		}
		limitCounter++
		if (limitCounter > limit) {break}
		var row, err = db.Query("select title, startAt, endAt from ad_data where id = ?", id)
		if err != nil {panic(err)}
		if (!row.Next()) {continue}
		var title string
		var startAt int64
		var endAt int64
		row.Scan(&title, &startAt, &endAt)
		if (startAt > now) {continue}
		if (endAt < now) {continue}			
		result = append(result, AdData{title, utils.TimestampToDateString(endAt)})
	}

	return result
}

func getNowTimestamp() int64 {
	return time.Now().UnixMilli()	
}
