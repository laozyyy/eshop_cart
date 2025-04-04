package service

import (
	"context"
	"eshop_cart/cache"
	"eshop_cart/database"
	"eshop_cart/log"
	"eshop_cart/model"
	"eshop_cart/util"
	"github.com/bytedance/sonic"
	"strconv"
)

func UpdateCart(uid string) {
	ctx := context.Background()
	result, err := cache.Client.HGetAll(ctx, util.GetKey(uid)).Result()
	result2, err := cache.Client.HGetAll(ctx, util.GetKeySelect(uid)).Result()
	if err != nil {
		log.Errorf("err: %v", err)
	}
	var items []model.CartItem
	for sku, q := range result {
		num, _ := strconv.ParseInt(q, 10, 32)
		item := model.CartItem{
			Sku:      sku,
			Quantity: int32(num),
			Selected: result2[sku] == "true",
		}
		items = append(items, item)
	}

	marshal, err := sonic.Marshal(items)
	err = database.InsertOrUpdateCart(nil, marshal, uid)
	if err != nil {
		log.Errorf("err: %v", err)
	}
}
