package database

import (
	"eshop_cart/log"
	"eshop_cart/model"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
	"time"
)

func getDBInstance(db *gorm.DB) *gorm.DB {
	if db == nil {
		if DB == nil {
			Init() // 初始化全局 DB
		}
		return DB
	}
	return db
}
func GetCartByUid(db *gorm.DB, uid string) ([]*model.CartItem, error) {
	db = getDBInstance(db)
	var ret []*model.CartItem
	var cart model.Cart
	err := db.Table("cart").Where("uid = ? AND is_deleted = 0", uid).Find(&cart).Error
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}
	if cart.CartData == nil {
		log.Info("购物车数据未持久化")
		return nil, err
	}
	err = sonic.Unmarshal(cart.CartData, &ret)
	if err != nil {
		log.Errorf("unmarshal error: %v", err)
		return nil, err
	}
	return ret, nil
}

func InsertOrUpdateCart(db *gorm.DB, cart []byte, uid string) error {
	db = getDBInstance(db)
	//cartModel := model.Cart{
	//	Uid:       uid,
	//	CartData:  cart,
	//	CreatedAt: time.Now(),
	//	UpdatedAt: time.Now(),
	//	IsDeleted: false,
	//}
	sql := `INSERT INTO cart (uid,cart_data,created_time,updated_time,is_deleted) VALUES (?,?,?,?,?)on duplicate key update cart_data = ?`
	err := db.Exec(sql, uid, cart, time.Now(), time.Now(), false, cart).Error
	//err := db.Table("cart").Save(&cartModel).Error
	if err != nil {
		log.Errorf("err: %v", err)
	}
	return err
}
