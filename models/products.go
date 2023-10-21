package models

import "pheet-fiber-backend/helper"

// *Entity เพื่อจะส่งข้อมูลออกไป *//
type Products struct {
	TableName   struct{}          `db:"products" json:"-" pk:"ID"`
	ID          string            `db:"id" json:"id" type:"string"`
	Title       string            `db:"title" json:"title" type:"string"`
	Description string            `db:"description" json:"description" type:"string"`
	Price       float64           `db:"price" json:"price" type:"float64"`
	CreatedAt   *helper.Timestamp `db:"created_at" json:"created_at" type:"timestamp"`
	UpdatedAt   *helper.Timestamp `db:"updated_at" json:"updated_at" type:"timestamp"`
	Categories  *Categories       `db:"-" json:"categories"`
	// ProductsCat []*ProductsCategories `json:"products_categories" db:"-" fk:"fk_field1:ID, fk_field2:ProductId"`
	Images []*Image `json:"images" db:"-" fk:"fk_field1:ID, fk_field2:ProductId"`
}
