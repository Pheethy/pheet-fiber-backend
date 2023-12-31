package models

import (
	"time"

	"github.com/Pheethy/psql/helper"
)

type Order struct {
	TableName struct{}          `json:"-" db:"orders" pk:"Id"`
	Id        string            `json:"id" db:"id" type:"string"`
	UserId    string            `json:"user_id" db:"user_id" type:"string"`
	Address   string            `json:"address" db:"address" type:"string"`
	Contact   string            `json:"contact" db:"contact" type:"string"`
	Status    string            `json:"status" db:"status" type:"string"`
	TotalPaid float64           `json:"total_paid" db:"-"` /* สำหรับคำนวณเงินในตระกร้า */
	CreatedAt *helper.Timestamp `json:"created_at" db:"created_at" type:"timestamp"`
	UpdatedAt *helper.Timestamp `json:"updated_at" db:"updated_at" type:"timestamp"`

	TransferSlip   *TransferSlip   `json:"transfer_slip" db:"-" fk:"fk_field1:Id, fk_field2:OrderId"`
	ProductsOrders []*ProductOrder `json:"products_orders" fk:"fk_field1:Id,fk_field2:OrderId"`
}

type TransferSlip struct {
	TableName struct{}          `json:"-" db:"transfer_slip" pk:"Id"`
	Id        string            `json:"id" db:"id" type:"string"`
	OrderId   string            `json:"order_id" db:"order_id" type:"string"`
	FileName  string            `json:"filename" db:"filename" type:"string"`
	Url       string            `json:"url" db:"url" type:"string"`
	CreatedAt *helper.Timestamp `json:"created_at" db:"created_at" type:"timestamp"`
}

type ProductOrder struct {
	TableName struct{} `json:"-" db:"products_orders" pk:"Id"`
	Id        string   `json:"id" db:"id" type:"string"`
	ProductId string   `json:"product_id" db:"product_id" type:"string"`
	OrderId   string   `json:"order_id" db:"order_id" type:"string"`
	Qty       int64    `json:"qty" db:"qty" type:"int64"`

	Products *Products `json:"products" db:"-" fk:"relation:one,fk_field1:ProductId,fk_field2:ID"`
}

func (p *Order) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.CreatedAt = &time
}

func (o *Order) FindingTotalPaid() {
	if len(o.ProductsOrders) > 0 {
		var totalPaid float64
		for index := range o.ProductsOrders {
			result := o.ProductsOrders[index].Products.Price * float64(o.ProductsOrders[index].Qty)
			totalPaid += result
		}
		o.TotalPaid = totalPaid
	}
}

func (p *Order) SetUpdatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.UpdatedAt = &time
}

func (p *TransferSlip) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.CreatedAt = &time
}
