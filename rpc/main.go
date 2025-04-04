package rpc

import (
	"context"
	"eshop_cart/kitex_gen/eshop/home"
	"eshop_cart/kitex_gen/eshop/home/goodsservice"
	"eshop_cart/log"
	"github.com/cloudwego/kitex/client"
	"strconv"
)

var goodsClient goodsservice.Client

func init() {
	var err error
	//goodsClient, err = goodsservice.NewClient("hello", client.WithHostPorts("localhost:8888"))
	goodsClient, err = goodsservice.NewClient("hello", client.WithHostPorts("117.72.72.114:20001"))
	if err != nil {
		log.Errorf("error: %v", err)
	}
}

func GetPrice(ctx context.Context, sku string) (float64, error) {
	req := &home.GetPriceRequest{
		Sku: sku,
	}
	resp, err := goodsClient.GetPrice(ctx, req)
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	if resp == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(resp, 64)
	return f, nil
}
