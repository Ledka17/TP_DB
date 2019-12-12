package model

import "time"

type Thread struct {
	Author		string
	Created		time.Time
	Forum		string
	Id			int32
	Message		string
	Slug		string
	Title		string
	Votes		int32
}

type ThreadUpdate struct {
	Message	string
	Title	string
}
