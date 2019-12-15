package model

type Vote struct {
	Nickname	string	`db:"nickname" json:"nickname"`
	Voice		int32	`db:"voice" json:"voice"`
}
