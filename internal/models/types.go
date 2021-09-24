package models

type User struct {
	Uid string `json:"uid"`
	Login string `json:"login"`
	Password string `json:"password"`
	Balance int `json:"balance"`
}

