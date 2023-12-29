package models

import (
	"github.com/Pheethy/psql/helper"
	"strings"
	"time"

	"github.com/gofrs/uuid"
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
	CategoriesId int64             `json:"categories_id" form:"categories_id" db:"-" type:"int64"` /* categories_id สำหรับการสร้าง products_categories */

	Categories *Categories `json:"categories" db:"-" fk:"fk_field1:ID, fk_field2:ProductId"` /* สำหรับการ Fetch Category มา Fill เพื่อดูว่า Product อยู่ Categories ไหน */
	Images     []*Image    `json:"images" db:"-" fk:"fk_field1:ID, fk_field2:ProductId"`
}

func (p *Products) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.CreatedAt = &time
}

func (p *Products) SetUpdatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.UpdatedAt = &time
}

func (p *Products) MergeProduct(exist *Products) {
	switch {
	case p.ID == "":
		p.ID = exist.ID
		fallthrough
	case p.Title == "":
		p.Title = exist.Title
		fallthrough
	case p.Description == "":
		p.Description = exist.Description
		fallthrough
	case p.Price == 0:
		p.Price = exist.Price
		fallthrough
	case p.CategoriesId == 0:
		p.CategoriesId = exist.Categories.Id
	}
	if len(p.Images) == 0 {
		p.Images = exist.Images
	}
}

func (p *Products) FindDeleteImage(exist *Products) ([]*uuid.UUID, []string) {
	var delIds = make([]*uuid.UUID, 0)
	var delURL = make([]string, 0)

	// Create a map for faster lookup of existing image UUIDs
	existImageMap := make(map[uuid.UUID]struct{})
	if len(p.Images) >= 1 { //เป็นการเช็คว่า product ต้องมีรูปอย่างน้อย 1 รูป
		for _, u := range p.Images {
			existImageMap[*u.ID] = struct{}{}
		}
	}

	// Iterate through new images and check if they exist in the existing map
	for _, u := range exist.Images {
		if _, exists := existImageMap[*u.ID]; !exists {
			delIds = append(delIds, u.ID) // Add the UUID to the delete list// Prefix to remove
			prefix := "https://storage.googleapis.com/pheethy-dev-bucket/"
			// Use strings.TrimLeft to remove the prefix
			result := strings.SplitAfter(u.Url, prefix)
			delURL = append(delURL, result[1])
		}
	}

	return delIds, delURL
}
