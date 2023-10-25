package models

import (
	"pheet-fiber-backend/helper"
	"time"
)

// *Entity เพื่อจะส่งข้อมูลออกไป *//
type Products struct {
	TableName    struct{}          `db:"products" json:"-" pk:"ID"`
	ID           string            `json:"id" form:"id" db:"id" type:"string"`
	Title        string            `json:"title" form:"title" db:"title" type:"string"`
	Description  string            `json:"description" form:"description" db:"description" type:"string"`
	Price        float64           `json:"price" form:"price" db:"price" type:"float64"`
	CreatedAt    *helper.Timestamp `json:"created_at" db:"created_at" type:"timestamp"`
	UpdatedAt    *helper.Timestamp `json:"updated_at" db:"updated_at" type:"timestamp"`
	CategoriesId int64             `json:"-" form:"categories_id" db:"-" type:"int64"` /* categories_id สำหรับการสร้าง products_categories */
	
	Categories   *Categories       `json:"categories" db:"-"`  /* สำหรับการ Fetch Category มา Fill เพื่อดูว่า Product อยู่ Categories ไหน */
	Images []*Image `json:"images" db:"-" fk:"fk_field1:ID, fk_field2:ProductId"`
}

func (p *Products) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.CreatedAt = &time
}

func (p *Products) SetUpdatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.UpdatedAt = &time
}