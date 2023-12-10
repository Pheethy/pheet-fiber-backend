package models

import (
	"github.com/Pheethy/psql/helper"
	"time"
)

type Order struct {
	Table        struct{}          `json:"-" db:"orders" pk:"Id"`
	Id           string            `json:"id" db:"id"`
	UserId       string            `json:"user_id" db:"user_id"`
	TransferSlip *TransferSlip     `json:"transfer_slip" db:"transfer_slip"`
	Products     []*ProductOrder   `json:"products"`
	Address      string            `json:"address" db:"address"`
	Contact      string            `json:"contact" db:"contact"`
	Status       string            `json:"status" db:"status"`
	TotalPaid    float64           `json:"total_paid" db:"-"` /* สำหรับคำนวณเงินในตระกร้า */
	CreatedAt    *helper.Timestamp `json:"created_at" db:"created_at"`
	UpdatedAt    *helper.Timestamp `json:"updated_at" db:"updated_at"`
}

type TransferSlip struct {
	Id        string            `json:"id" db:"id"`
	FileName  string            `json:"file_name" db:"file_name"`
	Url       string            `json:"url" db:"url"`
	CreatedAt *helper.Timestamp `json:"created_at" db:"created_at"`
}

type ProductOrder struct {
	Table   struct{}  `json:"-" db:"products_orders" pk:"Id"`
	Id      string    `json:"id" db:"id"`
	OrderId string    `json:"order_id" db:"order_id"`
	Qty     int       `json:"qty" db:"qty"`
	Product *Products `json:"product" db:"-"`
}

func (p *Order) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.CreatedAt = &time
}

func (p *Order) SetUpdatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.UpdatedAt = &time
}

func (p *TransferSlip) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	p.CreatedAt = &time
}