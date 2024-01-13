package models

import "github.com/gofrs/uuid"

type Categories struct {
	TableName struct{} `json:"-" db:"categories" pk:"Id"`
	Id        int64      `json:"id" db:"id" type:"int64"`
	Title     string   `json:"title" db:"title" type:"string"`

	ProductId  string     `json:"-" db:"product_id" type:"string"`
}

type ProductsCategories struct {
	TableName  struct{}   `json:"-" db:"products_categories" pk:"Id"`
	Id         *uuid.UUID `json:"id" db:"id" type:"uuid"`
	ProductId  string     `json:"product_id" db:"product_id" type:"string"`
	CategoryId int64        `json:"category_id" db:"category_id" type:"int64"`

	Category *Categories `json:"category" db:"-" fk:"fk_field1:CategoryId,fk_field2:Id"`
}
