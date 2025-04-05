package main

import (
	"context"
	"database/sql"
	"errors"
	"eshop_cart/cache"
	"eshop_cart/database"
	"eshop_cart/kitex_gen/eshop/cart"
	"eshop_cart/log"
	"eshop_cart/mq"
	"eshop_cart/rpc"
	"eshop_cart/util"
	"strconv"
	"time"
)

type CartServiceImpl struct{}

func (c CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemRequest) (r *cart.BaseResponse, err error) {
	uid := req.Uid
	sku := req.SkuId
	quantity := req.Quantity
	key := util.GetKey(uid)
	keySelect := util.GetKeySelect(uid)
	exists, err := isCacheExists(ctx, key)
	if exists {
		// 2. 缓存存在则更新数量
		err := cache.Client.HIncrBy(ctx, key, sku, int64(quantity)).Err()
		if err != nil {
			log.Errorf("获取缓存错误，error: %v", err)
			errStr := "服务器内部错误"
			return &cart.BaseResponse{
				Code:   500,
				ErrStr: &errStr,
			}, err
		}
		// 更新过期时间
		cache.Client.Expire(ctx, key, 24*time.Hour)
		cache.Client.Expire(ctx, keySelect, 24*time.Hour)
		go mq.SendSaveCartMessage(uid)
		return &cart.BaseResponse{
			Code:   200,
			ErrStr: nil,
		}, nil
	}

	//3. 查询并回填缓存
	ok, err := GetFromDBAndCacheCart(ctx, uid)
	if err != nil {
		log.Errorf("内部错误，err: %v", err)
		errStr := "服务器内部错误"
		return &cart.BaseResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}
	if ok {
		// 修改添加的sku数量
		err = cache.Client.HIncrBy(ctx, key, sku, int64(quantity)).Err()
	} else {
		// 4. 若数据库不存在，插入新商品到缓存
		_ = cache.Client.HIncrBy(ctx, key, sku, int64(quantity)).Err()
		_ = cache.Client.HSet(ctx, keySelect, sku, false).Err()
		if err != nil {
			log.Errorf("内部错误，err: %v", err)
			errStr := "服务器内部错误"
			return &cart.BaseResponse{
				Code:   500,
				ErrStr: &errStr,
			}, err
		}

	}
	// 6. 异步更新数据库
	go mq.SendSaveCartMessage(uid)
	return &cart.BaseResponse{
		Code:   200,
		ErrStr: nil,
	}, nil
}

func GetFromDBAndCacheCart(ctx context.Context, uid string) (bool, error) {
	exists, _ := isCacheExists(ctx, uid)
	if exists {
		return true, nil
	}
	// 不在缓存，查询数据库
	key := util.GetKey(uid)
	keySelect := util.GetKeySelect(uid)
	keyPrice := util.GetKeyPrice(uid)
	// 删除可能存在的的旧总价缓存
	cache.Client.Del(ctx, keyPrice)
	cartItems, err := database.GetCartByUid(nil, uid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("数据库错误，err: %v", err)
		return false, err
	}
	// 4. 数据库存在则回填缓存
	if cartItems != nil {
		for _, item := range cartItems {
			var err error
			err = cache.Client.HSet(ctx, key, item.Sku, item.Quantity).Err()
			if err != nil {
				log.Errorf("缓存回填失败: %v", err)
				return false, err
			}
			_ = cache.Client.HSet(ctx, keySelect, item.Sku, item.Selected).Err()
		}
		cache.Client.Expire(ctx, key, 24*time.Hour)
		cache.Client.Expire(ctx, keySelect, 24*time.Hour)
		return true, nil
	}
	return false, nil
}

func (c CartServiceImpl) GetList(ctx context.Context, req *cart.PageRequest) (r *cart.PageResponse, err error) {
	key := util.GetKey(req.Uid)
	keySelect := util.GetKeySelect(req.Uid)
	//3. 查询并回填缓存
	ok, err := GetFromDBAndCacheCart(ctx, req.Uid)
	if err != nil {
		log.Errorf("内部错误，err: %v", err)
		errStr := "服务器内部错误"
		return &cart.PageResponse{
			Info:  &errStr,
			Price: "",
		}, err
	}
	if ok {
		// todo 分页逻辑
		skus, err := cache.Client.HGetAll(ctx, key).Result()
		skusSelect, err := cache.Client.HGetAll(ctx, keySelect).Result()
		if err != nil {
			log.Errorf("err: %v", err)
			errStr := "服务器内部错误"
			return &cart.PageResponse{
				Info: &errStr,
			}, err
		}
		var items []*cart.CartItem
		for sku, q := range skus {
			num, _ := strconv.ParseInt(q, 10, 32)
			item := cart.CartItem{
				Sku:      sku,
				Quantity: int32(num),
				Selected: skusSelect[sku] == "1",
			}
			items = append(items, &item)
		}
		price := computePrice(ctx, items)
		var str = "success"
		return &cart.PageResponse{
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
			IsEnd:    false,
			Items:    items,
			Info:     &str,
			Price:    strconv.FormatFloat(price, 'f', 2, 64),
		}, nil
	}
	str := "购物车无数据"
	return &cart.PageResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    false,
		Items:    nil,
		Info:     &str,
	}, nil
}

func (c CartServiceImpl) UpdateItem(ctx context.Context, req *cart.UpdateRequest) (r *cart.UpdateResponse, err error) {
	key := util.GetKey(req.Uid)
	var price float64 = 0
	keySelect := util.GetKeySelect(req.Uid)
	exists, err := isCacheExists(ctx, key)
	keyPrice := util.GetKeyPrice(req.Uid)
	if exists {
		exists, err = isCacheExists(ctx, keyPrice)
		// lastSelected 是这次要修改的 sku 之前的选中状态
		lastSelected, _ := cache.Client.HGet(ctx, keySelect, req.SkuId).Result()
		cache.Client.HSet(ctx, keySelect, req.SkuId, req.Selected)
		// 如果存在 cart 和 cart_price
		if exists {
			tmp, _ := cache.Client.Get(ctx, keyPrice).Result()
			price, _ = strconv.ParseFloat(tmp, 64)
			if lastSelected == "0" && *req.Selected {
				getPrice, err := rpc.GetPrice(ctx, req.SkuId)
				price += getPrice * float64(*req.Quantity)
				if err != nil {
					log.Errorf("内部错误，err: %v", err)
				}
				log.Infof("fianl price: %f, quantity: %d", price, req.Quantity)
			} else if lastSelected == "1" && !*req.Selected {
				getPrice, err := rpc.GetPrice(ctx, req.SkuId)
				price -= getPrice * float64(*req.Quantity)
				if err != nil {
					log.Errorf("内部错误，err: %v", err)
				}
				log.Infof("fianl price: %f, quantity: %d", price, req.Quantity)
			} else if lastSelected == "1" && *req.Selected {
				// 只有数量改变的情况
				getPrice, err := rpc.GetPrice(ctx, req.SkuId)
				if err != nil {
					log.Errorf("内部错误，err: %v", err)
				}
				lastQuantity, _ := cache.Client.HGet(ctx, key, req.SkuId).Result()
				lastQuantityInt, err := strconv.ParseInt(lastQuantity, 10, 64)
				price += getPrice * float64(int64(*req.Quantity)-lastQuantityInt)
				log.Infof("fianl price: %f, quantity: %d, lastQuantity: %d", price, req.Quantity, lastQuantityInt)
			}
			_ = cache.Client.Set(ctx, keyPrice, price, time.Hour*24).Err()
		} else {
			skusSelect, _ := cache.Client.HGetAll(ctx, keySelect).Result()
			skus, _ := cache.Client.HGetAll(ctx, key).Result()
			var items []*cart.CartItem
			for sku, q := range skus {
				num, _ := strconv.ParseInt(q, 10, 32)
				item := cart.CartItem{
					Sku:      sku,
					Quantity: int32(num),
					Selected: skusSelect[sku] == "1",
				}
				items = append(items, &item)
			}
			price = computePrice(ctx, items)

			_ = cache.Client.Set(ctx, keyPrice, price, time.Hour*24).Err()
		}
	} else {
		// 查询并回填缓存
		ok, err := GetFromDBAndCacheCart(ctx, req.Uid)
		if err != nil || !ok {
			log.Errorf("内部错误，err: %v", err)
			errStr := "服务器内部错误"
			return &cart.UpdateResponse{
				Code:   500,
				ErrStr: &errStr,
			}, err
		}
		skus, _ := cache.Client.HGetAll(ctx, keySelect).Result()
		skusWithQuantity, _ := cache.Client.HGetAll(ctx, key).Result()
		for sku, selected := range skus {
			q, _ := strconv.Atoi(skusWithQuantity[sku])
			if sku == req.SkuId && *req.Selected {
				getPrice, err := rpc.GetPrice(ctx, sku)
				if err != nil {
					log.Errorf("内部错误，err: %v", err)
				}
				price += getPrice * float64(q)
			} else if selected == "true" {
				getPrice, err := rpc.GetPrice(ctx, sku)
				if err != nil {
					log.Errorf("内部错误，err: %v", err)
				}
				price += getPrice * float64(q)
			}
		}
		_ = cache.Client.Set(ctx, keyPrice, price, time.Hour*24).Err()
	}
	return &cart.UpdateResponse{
		Price:  strconv.FormatFloat(price, 'f', 2, 64),
		Code:   0,
		ErrStr: nil,
	}, nil
	//panic("")
}

func computePrice(ctx context.Context, items []*cart.CartItem) float64 {
	var price float64
	for _, item := range items {
		if item.Selected {
			getPrice, err := rpc.GetPrice(ctx, item.Sku)
			if err != nil {
				log.Errorf("内部错误，err: %v", err)
			}
			price += getPrice * float64(item.Quantity)
		}
	}
	return price
}

func isCacheExists(ctx context.Context, key string) (bool, error) {
	exists, err := cache.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("获取缓存错误，error: %v", err)
		return false, err
	}
	return exists == 1, nil
}

func (c CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteRequest) (r *cart.BaseResponse, err error) {
	key := util.GetKey(req.Uid)
	keySelect := util.GetKeySelect(req.Uid)
	keyPrice := util.GetKeyPrice(req.Uid)
	ok, err := GetFromDBAndCacheCart(ctx, req.Uid)
	if err != nil {
		log.Errorf("内部错误，err: %v", err)
		errStr := "内部错误"
		return &cart.BaseResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}
	if !ok {
		log.Errorf("意外的删除，err: %v", err)
		errStr := "意外的删除"
		return &cart.BaseResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}
	// 删除数据库，然后删除两个hash中的对应sku的内容，然后删掉价格
	err = database.DeleteAllCart(nil, req.Uid)
	cache.Client.HDel(ctx, key, req.Skus...)
	cache.Client.HDel(ctx, keySelect, req.Skus...)
	cache.Client.Del(ctx, keyPrice)
	// 最后发送mq写库
	go mq.SendSaveCartMessage(req.Uid)
	return &cart.BaseResponse{
		Code:   200,
		ErrStr: nil,
	}, nil
}
