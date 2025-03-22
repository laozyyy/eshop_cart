package main

import (
	"context"
	"eshop_cart/kitex_gen/eshop/cart"
	"eshop_cart/log"
	"eshop_cart/service"
	"fmt"
	"testing"
)

func TestCartServiceImpl_AddItem(t *testing.T) {
	type args struct {
		ctx context.Context
		req *cart.AddItemRequest
	}
	tests := []struct {
		args args
	}{
		{
			args: args{
				ctx: context.Background(),
				req: &cart.AddItemRequest{
					SkuId:    "12345",
					Quantity: 1342,
					Uid:      "930fbafc-f9ed-458f-a1cf-768d65f8825e",
				},
			},
		},
		{
			args: args{
				ctx: context.Background(),
				req: &cart.AddItemRequest{
					SkuId:    "12345d",
					Quantity: 13421,
					Uid:      "930fbafc-f9ed-458f-a1cf-768d65f8825e",
				},
			},
		},
		{
			args: args{
				ctx: context.Background(),
				req: &cart.AddItemRequest{
					SkuId:    "1234345",
					Quantity: 1,
					Uid:      "930fbafc-f9ed-458f-a1cf-768d65f8825e",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("1", func(t *testing.T) {
			c := CartServiceImpl{}
			_, err := c.AddItem(tt.args.ctx, tt.args.req)
			if err != nil {
				log.Errorf("err: %v", err)
			}

		})
	}
	service.UpdateCart()
}

func TestCartServiceImpl_GetList(t *testing.T) {
	type args struct {
		ctx context.Context
		req *cart.PageRequest
	}
	tests := []struct {
		args args
	}{
		{
			args: args{
				ctx: context.Background(),
				req: &cart.PageRequest{
					PageSize: 1,
					PageNum:  1,
					Uid:      "930fbafc-f9ed-458f-a1cf-768d65f8825e",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("1", func(t *testing.T) {
			c := CartServiceImpl{}
			res, err := c.GetList(tt.args.ctx, tt.args.req)
			if err != nil {
				log.Errorf("err: %v", err)
			}
			items := res.Items
			for _, item := range items {
				fmt.Printf("%+v", *item)
			}
		})
	}
	service.UpdateCart()
}
