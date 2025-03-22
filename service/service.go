package service

import (
	"context"
	"eshop_cart/cache"
	"eshop_cart/database"
	"eshop_cart/log"
	"eshop_cart/model"
	"fmt"
	"github.com/bytedance/sonic"
	"strconv"
)

func UpdateCart() {
	ctx := context.Background()
	uid := "930fbafc-f9ed-458f-a1cf-768d65f8825e"
	key := fmt.Sprintf("cart:{%s}", uid)
	result, err := cache.Client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Errorf("err: %v", err)
	}
	var items []model.CartItem
	for sku, q := range result {
		num, _ := strconv.ParseInt(q, 10, 32)
		item := model.CartItem{
			Sku:      sku,
			Quantity: int32(num),
		}
		items = append(items, item)
	}

	marshal, err := sonic.Marshal(items)
	err = database.InsertOrUpdateCart(nil, marshal, uid)
	if err != nil {
		log.Errorf("err: %v", err)
	}
}
