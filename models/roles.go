package models

type Roles struct {
	TableName struct{} `json:"-" db:"roles" pk:"Id"`
	Id        int64    `json:"id" db:"id" type:"int64"`
	Title     string   `json:"title" db:"title" type:"string"`
}
