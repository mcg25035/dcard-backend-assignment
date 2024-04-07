package dataHandler

import (
	"github.com/mcg25035/assignment/conditions"
)

type AdCondition struct {
	AgeStart int `json:"ageStart"`
	AgeEnd int `json:"ageEnd"`
	Gender []conditions.Gender `json:"gender"`
	Country []conditions.Country `json:"country"`
	Platform []conditions.Platform `json:"platform"`
}

type AdConditionRaw struct {
	AgeStart  int      `json:"ageStart"`
	AgeEnd    int      `json:"ageEnd"`
	Country   []conditions.Country `json:"country"`
	Platform  []conditions.Platform `json:"platform"`
	Gender  conditions.Gender `json:"gender"`
} 

type AdRequest struct {
	Title      string    `json:"title"`
	StartAt    string    `json:"startAt"` 
	EndAt      string    `json:"endAt"`   
	Conditions []AdConditionRaw `json:"conditions"`
}

type AdDataWithStartAt struct {
	Title string `json:"title"`
	EndAt string `json:"endAt"`
	StartAt string `json:"startAt"`
}

type AdData struct {
	Title string `json:"title"`
	EndAt string `json:"endAt"`
	StartAt string `json:"-"`
}
