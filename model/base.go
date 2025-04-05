package model

import (
	"eshop_cart/kitex_gen/eshop/cart"
	"time"
)

type CartItem struct {
	Sku      string
	Quantity int32
	Selected bool
}

type Cart struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	Uid       string    `gorm:"column:uid" json:"uid"`
	CartData  []byte    `gorm:"column:cart_data" json:"cart_data"`
	CreatedAt time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedAt time.Time `gorm:"column:updated_time" json:"updated_time"`
	IsDeleted bool      `gorm:"column:is_deleted" json:"is_deleted"`
}

type BySku []*cart.CartItem

// Len 返回数组的长度
func (b BySku) Len() int {
	return len(b)
}

// Less 比较两个元素的 Age 字段
func (b BySku) Less(i, j int) bool {
	return b[i].Sku < b[j].Sku
}

// Swap 交换两个元素的位置
func (b BySku) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
