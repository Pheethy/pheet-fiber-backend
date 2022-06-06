package models

//*Entity เพื่อจะส่งข้อมูลออกไป *//
type Product struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
	Price int `db:"price" json:"price"`
	Description string `db:"description" json:"description"`
}

type User struct {
	Id int `db:"id" json:"id"`
	UserName string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	NickName string `db:"nickname" json:"nickname"`
}