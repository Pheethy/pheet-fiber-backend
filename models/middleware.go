package models

type Role struct {
	Id int `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}