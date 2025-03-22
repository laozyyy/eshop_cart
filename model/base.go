package model

import "time"

type CartItem struct {
	Sku      string
	Quantity int32
}

type Cart struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	Uid       string    `gorm:"column:uid" json:"uid"`
	CartData  []byte    `gorm:"column:cart_data" json:"cart_data"`
	CreatedAt time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedAt time.Time `gorm:"column:updated_time" json:"updated_time"`
	IsDeleted bool      `gorm:"column:is_deleted" json:"is_deleted"`
}
