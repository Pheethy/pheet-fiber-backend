package models

type Catagory struct {
	Id    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}
