package repository

import (
	"database/sql"
	"fmt"
	"log"
)

func checkErr(err error) {
	if err != sql.ErrNoRows && err != nil {
		log.Panic(err)
	}
}

func getOrder(desc bool) string {
	if desc {
		return "desc"
	}
	return "asc"
}

func getFilterLimit (limit int) string {
	filterLimit := ""
	if limit > 0 {
		filterLimit = fmt.Sprintf(" limit %d ", limit)
	}
	return filterLimit
}

func getFilterId (order string, id int) string {
	if id == -1 {
		return ""
	}
	sign := getSign(order)
	filterId := fmt.Sprintf(" and id %s %d ", sign, id)
	return filterId
}

func getFilterSince(order string, since string) string{
	if since == "" {
		return ""
	}

	sign := getSign(order)
	filterSince := fmt.Sprintf(" and created %s= '%s' ", sign, since)
	return filterSince
}

func getFilterSinceByUserName(order string, since string) string{
	if since == "" {
		return ""
	}
	sign := getSign(order)
	filterSince := fmt.Sprintf(" and lower(u.nickname) %s lower('%s') ", sign, since)

	return filterSince
}

func getSign(order string) string {
	sign := ">"
	if order == "desc" {
		sign = "<"
	}
	return sign
}

func have(elem int64, array []int64) bool {
	for _, current := range array {
		if current == elem {
			return true
		}
	}
	return false
}
