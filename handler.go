package main

import (
	"context"
	"database/sql"
	"errors"
	"eshop_cart/cache"
	"eshop_cart/database"
	"eshop_cart/kitex_gen/eshop/cart"
	"eshop_cart/log"
	"fmt"
	"strconv"
	"time"
)

type CartServiceImpl struct{}

func (c CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemRequest) (r *cart.BaseResponse, err error) {
	uid := req.Uid
	sku := req.SkuId
	quantity := req.Quantity
	key := getKey(uid)
	exists, err := cache.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("获取缓存错误，error: %v", err)
		errStr := "服务器内部错误"
		return &cart.BaseResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}
	if exists == 1 {
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
		return &cart.BaseResponse{
			Code:   200,
			ErrStr: nil,
		}, nil
	}

	// 3. 不在缓存，查询数据库
	cartItems, err := database.GetCartByUid(nil, uid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("数据库错误，err: %v", err)
		errStr := "服务器内部错误"
		return &cart.BaseResponse{
			Code:   500,
			ErrStr: &errStr,
		}, nil
	}

	// 4. 数据库存在则回填缓存
	if cartItems != nil {
		for _, item := range cartItems {
			var err error
			if item.Sku == sku {
				err = cache.Client.HSet(ctx, key, item.Sku, item.Quantity+quantity).Err()
			} else {
				err = cache.Client.HSet(ctx, key, item.Sku, item.Quantity).Err()
			}
			if err != nil {
				log.Errorf("缓存回填失败: %v", err)
				errStr := "服务器内部错误"
				return &cart.BaseResponse{
					Code:   500,
					ErrStr: &errStr,
				}, nil
			}
		}
		cache.Client.Expire(ctx, key, 24*time.Hour)
		return &cart.BaseResponse{
			Code:   200,
			ErrStr: nil,
		}, nil
	}

	// 5. 数据库不存在，插入新商品到缓存
	if err := cache.Client.HIncrBy(ctx, key, sku, int64(quantity)).Err(); err != nil {
		log.Errorf("内部错误，err: %v", err)
		errStr := "服务器内部错误"
		return &cart.BaseResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}

	// 6. 异步更新数据库（todo 改为消息队列）
	return &cart.BaseResponse{
		Code:   200,
		ErrStr: nil,
	}, nil
}

func getKey(uid string) string {
	return fmt.Sprintf("cart:{%s}", uid)
}

func (c CartServiceImpl) GetList(ctx context.Context, req *cart.PageRequest) (r *cart.PageResponse, err error) {
	key := getKey(req.Uid)
	exists, err := cache.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("获取缓存错误，error: %v", err)
		errStr := "服务器内部错误"
		return &cart.PageResponse{
			Info: &errStr,
		}, err
	}
	if exists == 1 {
		// todo 分页逻辑
		result, err := cache.Client.HGetAll(ctx, key).Result()
		if err != nil {
			log.Errorf("err: %v", err)
			errStr := "服务器内部错误"
			return &cart.PageResponse{
				Info: &errStr,
			}, err
		}
		var items []*cart.CartItem
		for sku, q := range result {
			num, _ := strconv.ParseInt(q, 10, 32)
			item := cart.CartItem{
				Sku:      sku,
				Quantity: int32(num),
			}
			items = append(items, &item)
		}
		var str = "success"
		return &cart.PageResponse{
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
			IsEnd:    false,
			Items:    items,
			Info:     &str,
		}, nil
	}

	cartItems, err := database.GetCartByUid(nil, req.Uid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("数据库错误，err: %v", err)
		errStr := "服务器内部错误"
		return &cart.PageResponse{
			Info: &errStr,
		}, nil
	}
	// 缓存到redis
	for _, item := range cartItems {
		var err error
		err = cache.Client.HSet(ctx, key, item.Sku, item.Quantity).Err()
		if err != nil {
			log.Errorf("缓存回填失败: %v", err)
			errStr := "服务器内部错误"
			return &cart.PageResponse{
				Info: &errStr,
			}, nil
		}
	}
	cache.Client.Expire(ctx, key, 24*time.Hour)
	var items []*cart.CartItem
	for _, item := range cartItems {
		item := cart.CartItem{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		}
		items = append(items, &item)
	}
	var str = "success"
	return &cart.PageResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    false,
		Items:    items,
		Info:     &str,
	}, nil
}

func (c CartServiceImpl) UpdateItem(ctx context.Context, req *cart.UpdateRequest) (r *cart.UpdateResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (c CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteRequest) (r *cart.BaseResponse, err error) {
	//TODO implement me
	panic("implement me")
}
