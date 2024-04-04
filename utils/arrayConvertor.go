package utils

import (
	"strconv"
	"strings"
	"github.com/mcg25035/assignment/conditions"
	"math/rand"
)

func IntArrayToString(arr []int) string {
	str := ""
	var indx = 0
	for _, v := range arr {
		str += string(v)
		if indx < len(arr) - 1 {
			str += ","
		}
		indx++
	}
	return str
}

func StringToIntArray(str string) []int {
	var intRaw = strings.Split(str, ",")
	var result []int
	for _, v := range intRaw {
		var num, _ = strconv.Atoi(v)
		result = append(result, num)
	}
	return result
}

func CountryArrayToString(arr []conditions.Country) string {
	str := ""
	var indx = 0
	for _, v := range arr {
		str += string(v)
		if indx < len(arr) - 1 {
			str += ","
		}
		indx++
	}
	return str
}

func StringToCountryArray(str string) []conditions.Country {
	var countryRaw = strings.Split(str, ",")
	var result []conditions.Country
	for _, v := range countryRaw {
		result = append(result, conditions.Country(v))
	}
	return result
}

func GenderArrayToString(arr []conditions.Gender) string {
	str := ""
	var indx = 0
	for _, v := range arr {
		str += string(v)
		if indx < len(arr) - 1 {
			str += ","
		}
		indx++
	}
	return str
}

func StringToGenderArray(str string) []conditions.Gender {
	var genderRaw = strings.Split(str, ",")
	var result []conditions.Gender
	for _, v := range genderRaw {
		result = append(result, conditions.Gender(v))
	}
	return result
}

func PlatformArrayToString(arr []conditions.Platform) string {
	str := ""
	var indx = 0
	for _, v := range arr {
		str += string(v)
		if indx < len(arr) - 1 {
			str += ","
		}
		indx++
	}
	return str
}

func StringToPlatformArray(str string) []conditions.Platform {
	var platformRaw = strings.Split(str, ",")
	var result []conditions.Platform
	for _, v := range platformRaw {
		result = append(result, conditions.Platform(v))
	}
	return result
}

func GenerateRandomString() string {
	var result = "";
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
	for i := 0; i < 8; i++ {
		result += string(charset[rand.Intn(len(charset))])
	}
	return result
}